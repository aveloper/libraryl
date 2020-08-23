package db

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//Library stores basic parameters of Library Group
type Library struct {
	ID          string     `json:"id" bson:"_id"`
	LibraryName string     `json:"library_name" bson:"library_name"`
	GroupID     string     `json:"group_id" bson:"group_id"`
	URL         string     `json:"url" bson:"url"`
	Keywords    []string   `json:"keywords" bson:"keywords"`
	Artifacts   []Artifact `json:"artifacts" bson:"artifacts"`
}

//Stage stores the latest release in each stage
type Stage struct {
	Stable string `json:"stable" bson:"stable"`
	Alpha  string `json:"alpha" bson:"alpha"`
	Beta   string `json:"beta" bson:"beta"`
	RC     string `json:"rc" bson:"rc"`
}

//Artifact stores the dependencies of a Library Group
type Artifact struct {
	ArtifactName string `json:"artifact_name" bson:"artifact_name"`
	Stages       Stage  `json:"stages" bson:"stages"`
	Comments     string `json:"comments,omitempty" bson:"comments"`
	KotlinOnly   bool   `json:"kotlin_only,omitempty" bson:"kotlin_only"`
	Replaces     string `json:"replaces,omitempty" bson:"replaces"`
	Processor    bool   `json:"processor,omitempty" bson:"processor"`
	Testing      bool   `json:"testing,omitempty" bson:"testing"`
}

//Add a new library to db
func (l *Library) Add(db *mongo.Collection, ctx context.Context) error {
	if l.ID == "" {
		l.ID = uuid.New().String()
	}

	_, err := db.InsertOne(ctx, l)
	if err != nil {
		return err
	}

	return nil
}

//Update the library present in db
func (l *Library) Update(db *mongo.Collection, ctx context.Context) error {
	filter := struct {
		ID string `bson:"_id"`
	}{
		ID: l.ID,
	}

	_, err := db.UpdateOne(ctx, filter, l)
	if err != nil {
		return err
	}

	return nil
}

//Get the Library details
func (l *Library) Get(db *mongo.Collection, ctx context.Context) error {
	filter := struct {
		ID string `bson:"_id"`
	}{
		ID: l.ID,
	}

	r := db.FindOne(ctx, filter)
	if r.Err() != nil {
		return r.Err()
	}

	err := r.Decode(l)
	if err != nil {
		return err
	}

	return nil
}

//GetAll the Library details in db
func (l *Library) GetAll(db *mongo.Collection, ctx context.Context) ([]*Library, error) {
	libraries := make([]*Library, 0)

	c, err := db.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	for c.Next(ctx) {
		var lib *Library

		err = c.Decode(lib)
		if err != nil {
			return nil, err
		}

		libraries = append(libraries, lib)
	}

	return libraries, nil
}

//Delete the Library from db
func (l *Library) Delete(db *mongo.Collection, ctx context.Context) error {

	_, err := db.DeleteOne(ctx, bson.M{"_id": l.ID})
	if err != nil {
		return err
	}

	return nil
}
