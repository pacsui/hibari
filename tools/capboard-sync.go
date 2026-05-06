package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CapTransaction struct {
	MessageID string    `bson:"message_id"`
	ChannelID string    `bson:"channel_id"`
	GuildID   string    `bson:"guild_id"`
	GiverID   string    `bson:"giver_id"`
	TakerID   string    `bson:"taker_id"`
	CreatedAt time.Time `bson:"created_at"`
}

type UniqueMessage struct {
	MessageID string `bson:"_id"`
	ChannelID string `bson:"channel_id"`
	GuildID   string `bson:"guild_id"`
	TakerID   string `bson:"taker_id"`
}

func main() {
	log.Info("Starting Capboard Synchronization Script...")
	mongoURI := os.Getenv("MONGODB_URI")
	discordToken := os.Getenv("DISCORD_TOKEN")

	if mongoURI == "" || discordToken == "" {
		log.Fatal("MONGODB_URI and DISCORD_TOKEN environment variables are required.")
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("MongoDB Connection Error: %v", err)
	}
	defer client.Disconnect(context.Background())
	capboard := client.Database("Hibari").Collection("Capboard")

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalf("Discord Session Error: %v", err)
	}

	ctx := context.Background()

	log.Info("Fetching unique messages from Database...")
	var uniqueMsgs []UniqueMessage
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$message_id"},
			{Key: "channel_id", Value: bson.D{{Key: "$first", Value: "$channel_id"}}},
			{Key: "guild_id", Value: bson.D{{Key: "$first", Value: "$guild_id"}}},
			{Key: "taker_id", Value: bson.D{{Key: "$first", Value: "$taker_id"}}},
		}}},
	}

	cursor, err := capboard.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatalf("Aggregation failed: %v", err)
	}
	if err = cursor.All(ctx, &uniqueMsgs); err != nil {
		log.Fatalf("Failed to decode unique messages: %v", err)
	}

	log.Infof("Found %d unique messages to verify.", len(uniqueMsgs))

	var statsAdded, statsRemoved, statsDeletedMsgs int

	for i, msg := range uniqueMsgs {
		log.Infof("[%d/%d] Checking Message ID: %s", i+1, len(uniqueMsgs), msg.MessageID)

		var existingCaps []CapTransaction
		capCursor, err := capboard.Find(ctx, bson.D{{Key: "message_id", Value: msg.MessageID}})
		if err == nil {
			capCursor.All(ctx, &existingCaps)
		}

		dbGivers := make(map[string]bool)
		for _, c := range existingCaps {
			dbGivers[c.GiverID] = true
		}
		discordGivers := make(map[string]bool)
		var lastID string
		messageExists := true

		for {
			users, err := dg.MessageReactions(msg.ChannelID, msg.MessageID, "🧢", 100, "", lastID)

			if err != nil {
				var restErr *discordgo.RESTError
				if errors.As(err, &restErr) && restErr.Response != nil && restErr.Response.StatusCode == 404 {
					messageExists = false
				} else {
					log.Errorf("Failed to fetch reactions for message %s: %v", msg.MessageID, err)
				}
				break
			}

			if len(users) == 0 {
				break
			}

			for _, u := range users {
				discordGivers[u.ID] = true
			}
			lastID = users[len(users)-1].ID
		}

		if !messageExists {
			res, _ := capboard.DeleteMany(ctx, bson.D{{Key: "message_id", Value: msg.MessageID}})
			log.Warnf("Message %s deleted from Discord. Removed %d orphaned cap records.", msg.MessageID, res.DeletedCount)
			statsDeletedMsgs++
			continue
		}

		for giverID := range discordGivers {
			if giverID == msg.TakerID {
				continue
			}
			if !dbGivers[giverID] {
				newCap := CapTransaction{
					MessageID: msg.MessageID,
					ChannelID: msg.ChannelID,
					GuildID:   msg.GuildID,
					GiverID:   giverID,
					TakerID:   msg.TakerID,
					CreatedAt: time.Now(),
				}
				_, err := capboard.InsertOne(ctx, newCap)
				if err != nil && !mongo.IsDuplicateKeyError(err) {
					log.Errorf("Failed to insert missing cap: %v", err)
				} else {
					log.Debugf("Added missing cap to DB! Giver: %s", giverID)
					statsAdded++
				}
			}
		}

		for giverID := range dbGivers {
			if !discordGivers[giverID] {
				filter := bson.D{
					{Key: "message_id", Value: msg.MessageID},
					{Key: "giver_id", Value: giverID},
				}
				_, err := capboard.DeleteOne(ctx, filter)
				if err != nil {
					log.Errorf("Failed to delete stale cap: %v", err)
				} else {
					log.Debugf("Removed un-reacted cap from DB! Giver: %s", giverID)
					statsRemoved++
				}
			}
		}

		time.Sleep(200 * time.Millisecond)
	}

	log.Info("=== Sync Complete! ===")
	log.Infof("Missing Caps Added to DB: %d", statsAdded)
	log.Infof("Removed Caps Deleted from DB: %d", statsRemoved)
	log.Infof("Deleted/Orphaned Messages Cleared: %d", statsDeletedMsgs)
}
