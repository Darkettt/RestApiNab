package main

import (
	//	"fmt"

	"fmt"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	"context"
	"log"

	//"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var databases, ctx = mongodb()

func main() {

	r := gin.Default()

	// r := gin.New()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())
	r.GET("/AllUsers", GetAllUsers)
	r.GET("/AllMusics", GetAllMusics)
	r.GET("/PlaylistByUser/:userId", GetPlaylistByUser)
	r.GET("/allMusicsByPlaylist/:playListId", GetAllMusicsByPlaylist)
	r.GET("/PlaylistByMusic/:musicId", GetPlaylistByMusic)

	r.Run(":8080")
}

func GetAllMusicsByPlaylist(c *gin.Context) {
	playLisId := c.Param("playListId")
	fmt.Println(playLisId)
	objectId, _ := primitive.ObjectIDFromHex(playLisId)
	databaseMusic, err := databases.Collection("musics").Find(ctx, bson.M{"playlistIds": objectId})
	if err != nil {
		panic(err)
	}

	var musicResult = []Music{}
	databaseMusic.All(ctx, &musicResult)
	c.JSON(http.StatusOK, musicResult)

}

func GetPlaylistByMusic(c *gin.Context) {
	musicid := c.Param("musicId")

	objectId, _ := primitive.ObjectIDFromHex(musicid)
	var musicResult = Music{}

	err := databases.Collection("musics").FindOne(ctx, bson.M{"_id": objectId}).Decode(&musicResult)
	if err != nil {
		panic(err)
	}
	var playlistsResult = []Playlist{}
	databasePlaylist, err := databases.Collection("playlists").Find(ctx, bson.M{"_id": bson.M{"$in": musicResult.PlaylistIds}})

	databasePlaylist.All(ctx, &playlistsResult)
	c.JSON(http.StatusOK, playlistsResult)
}

func GetPlaylistByUser(c *gin.Context) {
	userid := c.Param("userId")
	objectId, _ := primitive.ObjectIDFromHex(userid)
	databasePlaylist, err := databases.Collection("playlists").Find(ctx, bson.M{"userId": objectId})
	if err != nil {
		panic(err)
	}
	var playlistsResult = []Playlist{}

	databasePlaylist.All(ctx, &playlistsResult)
	fmt.Println(len(playlistsResult))

	c.JSON(http.StatusOK, playlistsResult)

}

func GetAllUsers(c *gin.Context) {
	databaseUser, err := databases.Collection("users").Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}

	var userResult = []User{}

	databaseUser.All(ctx, &userResult)

	c.JSON(http.StatusOK, userResult)
}

func GetAllMusics(c *gin.Context) {
	databaseMusic, err := databases.Collection("musics").Find(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	var musicResult = []Music{}

	databaseMusic.All(ctx, &musicResult)
	c.JSON(http.StatusOK, musicResult)
}

func mongodb() (*mongo.Database, context.Context) {

	mongodb := "mongodb+srv://admin:admin@cluster0.n5ow2vn.mongodb.net/streaming?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongodb))
	if err != nil {
		log.Fatal(err)
	}
	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases := client.Database("streaming")

	return databases, ctx
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Age       int64              `json:"age,omitempty" bson:"age,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type Music struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string               `json:"name,omitempty" bson:"name,omitempty"`
	Artist      string               `json:"artist,omitempty" bson:"artist,omitempty"`
	PlaylistIds []primitive.ObjectID `json:"playlistIds,omitempty" bson:"playlistIds,omitempty"`
	CreatedAt   time.Time            `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time            `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
type Playlist struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	userId    primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
