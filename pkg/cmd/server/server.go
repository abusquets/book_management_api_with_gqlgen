package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/urfave/cli/v2"

	"github.com/abusquets/book_management_api_with_gqlgen/graph"
)

const defaultPort = "8080"

var DB *bun.DB

func run(port string) {
	if port == "" {
		port = defaultPort
	}

	// Make a connection with the database
	err := connectToDatabase()
	if err != nil {
		log.Fatalf("Unable to connect to the database: %s", err.Error())
	}
	fmt.Println("Successfully connected to the database")

	resolver := &graph.Resolver{
		DB: DB,
	}

	es := graph.NewExecutableSchema(graph.Config{Resolvers: resolver})
	srv := handler.NewDefaultServer(es)

	// Handler for GraphQL playground
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// Handler for GraphQL queries and mutations
	http.Handle("/query", srv)

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectToDatabase() error {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading app.env file")
	}

	// Build the database connection string
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	// Open a connection to the PostgreSQL database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Initialize the bun.DB instance with the PostgreSQL dialect
	DB = bun.NewDB(sqldb, pgdialect.New())

	// Add a query hook for debugging purposes
	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return DB.Ping()
}

func RunGraphqlServerCommand() *cli.Command {
	// go run cmd/booking/booking.go start-server --port 9090

	myFlags := []cli.Flag{
		&cli.StringFlag{
			Name:     "port",
			Value:    "",
			Usage:    "Port to run the server on",
			Required: false,
		},
	}

	command := &cli.Command{
		Name:  "start-server",
		Usage: "Start the server for the GraphQL API",
		Flags: myFlags,
		Action: func(cCtx *cli.Context) error {
			port := cCtx.String("port")
			run(port)
			return nil
		},
	}
	return command
}
