package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/jedib0t/go-pretty/table"
)

type DBConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type ListCommand struct{}

type GetCommand struct {
	SecretName string `arg:"" required:"" help:"Name of the secret to retrieve."`
}

var CLI struct {
	List ListCommand `cmd:"" help:"List all secrets"`
	Get  GetCommand  `cmd:"" help:"Get a specific secret"`
}

func createSecretsManagerClient() *secretsmanager.Client {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile("localstack"),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return secretsmanager.NewFromConfig(cfg)
}

func (l *ListCommand) Run() error {
	secretsManagerClient := createSecretsManagerClient()
	resp, err := secretsManagerClient.ListSecrets(
		context.TODO(),
		&secretsmanager.ListSecretsInput{},
	)

	if err != nil {
		log.Fatalf("Failed to list secrets, %v", err)
	}

	var rowHeader = table.Row{"Name", "ARN"}
	tw := table.NewWriter()
	tw.AppendHeader(rowHeader)

	for _, secret := range resp.SecretList {
		row := table.Row{aws.ToString(secret.Name), aws.ToString(secret.ARN)}
		tw.AppendRow(row)
	}

	fmt.Println(tw.Render())
	return nil
}

func (l *GetCommand) Run() error {

	client := createSecretsManagerClient()

	resp, err := client.GetSecretValue(
		context.TODO(),
		&secretsmanager.GetSecretValueInput{
			SecretId: &l.SecretName,
		},
	)

	if err != nil {
		log.Fatalf("Failed to get secret, %v", err)
	}

	var dbConfig DBConfig

	json.Unmarshal([]byte(aws.ToString(resp.SecretString)), &dbConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	jdbcURL := fmt.Sprintf("jdbc:postgresql://%s:%d/?user=%s&password=%s", dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password)

	fmt.Printf(jdbcURL)
	return nil
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
