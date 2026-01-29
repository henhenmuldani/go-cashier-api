package main

import (
	"log"      // Logging package
	"net/http" // HTTP server package
	"os"       // Operating system functionality package
	"strings"  // String manipulation package

	"github.com/spf13/viper" // Configuration package

	"go-cashier-api/database" // Import database package
	"go-cashier-api/handler"  // Import handler package
	"go-cashier-api/repository"
	"go-cashier-api/service" // Import service package
)

type Config struct {
	Port   string `mapstructure:"PORT"`    // Server port
	DBConn string `mapstructure:"DB_CONN"` // Database connection string
}

func main() {
	viper.AutomaticEnv()                                   // read in environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // replace dots with underscores
	// Load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	// Map environment variables to Config struct
	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DBCONN"),
	}

	//Setup database connection
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories, services, and handlers for products
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// Initialize repositories, services, and handlers for categories
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	// Setup HTTP server and routes
	mux := http.NewServeMux()

	// Define routes and their handlers
	mux.HandleFunc("/health", handler.HealthHandler)
	mux.HandleFunc("/api/products", productHandler.HandleProducts)
	mux.HandleFunc("/api/products/", productHandler.HandleProducts)
	mux.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	mux.HandleFunc("/api/categories/", categoryHandler.HandleCategories)

	log.Println("Server running on :" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, mux))
}
