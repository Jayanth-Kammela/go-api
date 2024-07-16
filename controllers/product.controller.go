package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jayanth-Kammela/go-api/database"
	"github.com/Jayanth-Kammela/go-api/models"
)

// CreateProduct creates a new product
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Request Body: %+v", product)

	collection := database.GetCollection()

	result, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		response := map[string]interface{}{
			"error":   err.Error(),
			"message": "Failed to create product",
			"status":  http.StatusInternalServerError,
		}
		// w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Retrieve the inserted product from the result to get the generated _id
	insertedID := result.InsertedID.(primitive.ObjectID)
	insertedProduct := models.Product{
		Id:          insertedID,
		Image:       product.Image,
		Title:       product.Title,
		Description: product.Description,
		Price:       product.Price,
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"data":    insertedProduct,
		"message": "Product created successfully",
		"status":  http.StatusCreated,
	}
	json.NewEncoder(w).Encode(response)
}

// GetProduct retrieves a single product by ID
func GetProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid product ID")
		return
	}

	collection := database.GetCollection()

	var product models.Product
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"message": "Product not found",
			"status":  http.StatusNotFound,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	// fmt.Print(product)
	response := map[string]interface{}{
		"message": "Product retrieved successfully",
		"data":    product,
		"status":  http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

// GetProducts retrieves all products
func GetProducts(w http.ResponseWriter, r *http.Request) {

	collection := database.GetCollection()

	var products []models.Product
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var product models.Product
		err := cursor.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	response := map[string]interface{}{
		"message": "Products retrieved successfully",
		"data":    products,
		"status":  http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateProduct updates an existing product by ID
func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := database.GetCollection()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"image":       updatedProduct.Image,
			"title":       updatedProduct.Title,
			"description": updatedProduct.Description,
			"price":       updatedProduct.Price,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated product from the database
	var fetchedProduct models.Product
	err = collection.FindOne(context.Background(), filter).Decode(&fetchedProduct)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"message": "Product not found",
			"status":  http.StatusNotFound,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"message": "Product updated successfully",
		"data":    fetchedProduct,
		"status":  http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteProduct deletes a product by ID
func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	collection := database.GetCollection()

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"message": "Product not found",
			"status":  http.StatusNotFound,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// fmt.Fprintf(w, "Product deleted successfully")
	response := map[string]interface{}{
		"message": "Product deleted successfully",
		"status":  http.StatusOK,
	}
	json.NewEncoder(w).Encode(response)
}
