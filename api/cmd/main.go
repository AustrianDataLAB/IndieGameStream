package main

import (
	"api/apis"
	"api/controllers"
	"api/repositories"
	"api/scripts"
	"api/services"
	"context"
	"database/sql"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	streamv1 "indiegamestream.com/indiegamestream/api/stream/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func setupRouter(db *sql.DB, azClient *azblob.Client) *gin.Engine {
	//Setup Gin
	r := gin.Default()
	//Cors
	r.Use(CORSMiddleware())

	//Repositories
	gamesRepository := repositories.GameRepository(db)

	//Apis
	k8sApi := apis.K8sService(k8sClient())
	azureApi := apis.AzureService(azClient)

	//Services
	gamesService := services.GameService(gamesRepository, k8sApi, azureApi)
	authService := services.AuthService()

	//Controllers
	gamesController := controllers.GameController(gamesService)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	//Upload a game
	r.POST("/games", authService.Authorize, gamesController.UploadGame)
	//Get all uploaded games
	r.GET("/games", authService.Authorize, gamesController.GetAllGames)
	//Get a specific game by its id
	r.GET("/games/:id", authService.Authorize, gamesController.GetGameById)
	//Delete a specific game, identified by its id
	r.DELETE("/games/:id", authService.Authorize, gamesController.DeleteGameById)

	return r
}

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(fmt.Sprintf("Failed to load .env file: %s", err))
		log.Println("Environment variables will be used")
	}
}

func setupDatabase() *sql.DB {
	//Create database if it is not existing yet.
	//We might have to remove this if we use an azure database
	scripts.CreateDatabaseIfNotExists(os.Getenv("MYSQL_DATABASE"))
	//Connect to the database
	db := scripts.ConnectToDatabase()
	//Check if database is online
	err := db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	//Check if we have new migrations and apply them
	scripts.MigrateDatabase(db)
	return db
}

func setupManagedClustersClient() *armcontainerservice.ManagedClustersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	clientFactory, err := armcontainerservice.NewClientFactory(
		os.Getenv("AZURERM_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	return clientFactory.NewManagedClustersClient()
}

func getKubeConfig() armcontainerservice.ManagedClustersClientListClusterUserCredentialsResponse {
	managedClustersClient := setupManagedClustersClient()

	res, err := managedClustersClient.ListClusterUserCredentials(context.Background(),
		os.Getenv("AZURERM_RESOURCE_GROUP_NAME"),
		os.Getenv("AZURE_AKS_CLUSTER_NAME"),
		&armcontainerservice.ManagedClustersClientListClusterUserCredentialsOptions{ServerFqdn: nil, Format: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}

	return res
}

func k8sClient() client.Client {
	//Try to read k8s config directly from environment
	restConfig, err := config.GetConfig()
	if err != nil || restConfig == nil {
		//If it fails, try to get it from azure
		kubeConfig := getKubeConfig()
		if len(kubeConfig.Kubeconfigs) == 0 {
			log.Fatalf("The kubeconfig request was successful but it's response body is empty")
		}
		if len(kubeConfig.Kubeconfigs) > 1 {
			log.Println("WARNING: Multiple kube-config's have been found. The first one will be used.")
		}

		clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeConfig.Kubeconfigs[0].Value)
		if err != nil {
			log.Fatalf("failed to create client config: %v", err)
		}
		restConfig, err = clientConfig.ClientConfig()
		if err != nil {
			log.Fatalf("failed to load kube config: %v", err)
		}
	}

	scheme, err := createScheme()
	if err != nil {
		log.Fatal(err.Error())
	}

	k8sc, err := client.New(
		restConfig,
		client.Options{Scheme: scheme},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	return k8sc
}

func createScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	// register all built-in types
	err := clientgoscheme.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	// register game scheme
	err = streamv1.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	return scheme, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func setupAzureBlobContainer(azClient *azblob.Client) {

	containerName := os.Getenv("AZURE_CONTAINER_NAME")
	containerClient := azClient.ServiceClient().NewContainerClient(containerName)

	if containerClient != nil {
		log.Println(fmt.Sprintf("Azure blob container with name %s exists already", containerName))
	} else {
		_, err := azClient.CreateContainer(context.Background(), containerName, nil)
		if err != nil {
			log.Fatal("Creating Azure Blob Container failed")
		}
		log.Println(fmt.Sprintf("Created Azure container %s", containerName))
	}
}

func main() {
	//Load config file
	loadConfig()

	//Setup Azure
	azClient := setupAzureBlobClient()
	setupAzureBlobContainer(azClient)

	//Setup database
	db := setupDatabase()
	defer db.Close()

	//Set Gin-gonic to debug or release mode
	gin.SetMode(os.Getenv("GIN_MODE"))

	//Setup Routes
	r := setupRouter(db, azClient)

	// Listen and Server in 0.0.0.0:8080
	err := r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setupAzureBlobClient() *azblob.Client {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/", os.Getenv("AZURE_STORAGE_ACCOUNT"))

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Initializing DefaultAzureCredential failed")
	}

	client, err := azblob.NewClient(url, credential, nil)
	if err != nil {
		log.Fatal("Initializing Azure Blob Client failed")
	}

	return client
}
