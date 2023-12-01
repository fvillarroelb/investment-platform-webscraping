package adapter

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository es la implementación de QueryRepository
type Repository struct {
	db                        *Database
	collectionInvestment      *mongo.Collection
	collectionInvestmentGroup *mongo.Collection
}

// NewRepository crea una nueva instancia de Repository
func NewRepository(db *Database, collectionInvestment *mongo.Collection, collectionInvestmentGroup *mongo.Collection) *Repository {
	//clientInversiones :=
	return &Repository{
		db:                        db,
		collectionInvestment:      collectionInvestment,
		collectionInvestmentGroup: collectionInvestmentGroup,
	}
}

// createUniqueIndex crea un índice único en la colección
func (repo *Repository) CreateUniqueIndex() error {
	// Crea un índice único en los campos relevantes (Ejemplo: Descripcion y Monto)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "nemo", Value: 1},
			{Key: "amount", Value: 1},
			{Key: "date", Value: 1},
			{Key: "platform", Value: 1},
			// Agrega otros campos según sea necesario
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.collectionInvestment.Indexes().CreateOne(context.Background(), indexModel)
	return err
}
func (repo *Repository) CreateUniqueIndexGroup() error {
	// Crea un índice único en los campos relevantes (Ejemplo: Descripcion y Monto)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "nemo", Value: 1},
			{Key: "amount", Value: 1},
			{Key: "date", Value: 1},
			{Key: "platform", Value: 1},
			// Agrega otros campos según sea necesario
		},
		Options: options.Index().SetUnique(true),
	}

	_, err := repo.collectionInvestmentGroup.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (repo *Repository) QueryCountInvestment() (int64, error) {

	docCount, err := repo.collectionInvestment.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Count Investment Collection	: ", docCount)
	return docCount, err
}

func (repo *Repository) InsertManyInversiones(moves []interface{}) (*mongo.InsertManyResult, error) {

	result, err := repo.collectionInvestment.InsertMany(context.Background(), moves)
	if err != nil {
		fmt.Println("Error insert many:", err)
		return result, err
	}
	// Imprime los IDs de los documentos insertados
	fmt.Println("Documentos insertados correctamente. IDs:", result)
	return result, err
}

func (repo *Repository) InsertManyInversionesAgrupadas(moves []interface{}) (*mongo.InsertManyResult, error) {

	result, err := repo.collectionInvestmentGroup.InsertMany(context.Background(), moves)
	if err != nil {
		fmt.Println("Error insert many:", err)
		return result, err
	}
	// Imprime los IDs de los documentos insertados
	fmt.Println("Documentos Agrupados insertados  correctamente. IDs:", result)
	return result, err
}

func (repo *Repository) DropCollectionInversiones() error {

	err := repo.collectionInvestment.Drop(context.Background())
	if err != nil {
		fmt.Println("Error al Drop Collection Inversiones:", err)
		return err
	}
	fmt.Println("Drop Collection Inversiones:", err)
	return err
}

func (repo *Repository) DropCollectionInversionesGroup() error {

	err := repo.collectionInvestmentGroup.Drop(context.Background())
	if err != nil {
		fmt.Println("Error al Drop Collection Inversiones:", err)
		return err
	}
	fmt.Println("Drop Collection Inversiones:", err)
	return err
}
