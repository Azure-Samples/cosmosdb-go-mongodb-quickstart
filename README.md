---
page_type: sample
languages:
- go
products:
- azure
description: "Build a Go application with MongoDB"
urlFragment: "update-this-to-unique-url-stub"
---

# Quickstart: Connect a Go application to Azure Cosmos DB's API for MongoDB

The sample application is a command-line based `todo` management tool written in Go. Azure Cosmos DB's API for MongoDB is [compatible with the MongoDB wire protocol](https://docs.microsoft.com/azure/cosmos-db/mongodb-introduction#wire-protocol-compatibility), making it possible for any MongoDB client driver to connect to it. This application uses the [Go driver for MongoDB](https://github.com/mongodb/mongo-go-driver) in a way that is transparent to the application that the data is stored in an Azure Cosmos DB database.

## Prerequisites

- An Azure account with an active subscription. [Create one for free](https://azure.microsoft.com/free). Or [try Azure Cosmos DB for free](https://azure.microsoft.com/try/cosmosdb/) without an Azure subscription. You can also use the [Azure Cosmos DB Emulator](https://aka.ms/cosmosdb-emulator) with the connection string `.mongodb://localhost:C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw==@localhost:10255/admin?ssl=true`.
- [Go](https://golang.org/) installed on your computer, and a working knowledge of Go.
- [Git](https://git-scm.com/downloads).
- If you don't want to use Azure Cloud Shell, [Azure CLI 2.0+](/cli/azure/install-azure-cli).

## Setup

Clone the application

```bash
git clone https://github.com/Azure-Samples/cosmosdb-go-mongodb-quickstart
```

Change into the directory where you cloned the application and build it (using `go build`).

```bash
cd monogdb-go-quickstart
go build -o todo
```

To confirm that the application was built properly.

```bash
./todo --help
```

To configure the application, export the connection string, MongoDB database and collection names as environment variables. 

```bash
export MONGODB_CONNECTION_STRING="mongodb://<COSMOSDB_ACCOUNT_NAME>:<COSMOSDB_PASSWORD>@<COSMOSDB_ACCOUNT_NAME>.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&maxIdleTimeMS=120000&appName=@<COSMOSDB_ACCOUNT_NAME>@"
```

> [!NOTE] 
> The `ssl=true` option is important because of Cosmos DB requirements. For more information, see [Connection string requirements](connect-mongodb-account.md#connection-string-requirements).
>

For the `MONGODB_CONNECTION_STRING` environment variable, replace the placeholders for `<COSMOSDB_ACCOUNT_NAME>` and `<COSMOSDB_PASSWORD>`

1. `<COSMOSDB_ACCOUNT_NAME>`: The name of the Azure Cosmos DB account
2. `<COSMOSDB_PASSWORD>`: The Azure Cosmos DB database key

```bash
export MONGODB_DATABASE=todo-db
export MONGODB_COLLECTION=todos
```

You can choose your preferred values for `MONGODB_DATABASE` and `MONGODB_COLLECTION` or leave them as is.

## Running the sample

To create a `todo`

```bash
./todo --create "Create an Azure Cosmos DB database account"
```

If successful, you should see an output with the MongoDB `_id` of the newly created document:

```bash
added todo ObjectID("5e9fd6befd2f076d1f03bd8a")
```

Create another `todo`

```bash
./todo --create "Get the MongoDB connection string using the Azure CLI"
```

List all the `todo`s

```bash
./todo --list all
```

You should see the ones you just added in a tabular format as such

```bash
+----------------------------+--------------------------------+-----------+
|             ID             |          DESCRIPTION           |  STATUS   |
+----------------------------+--------------------------------+-----------+
| "5e9fd6b1bcd2fa6bd267d4c4" | Create an Azure Cosmos DB      | pending   |
|                            | database account               |           |
| "5e9fd6befd2f076d1f03bd8a" | Get the MongoDB connection     | pending   |
|                            | string using the Azure CLI     |           |
+----------------------------+--------------------------------+-----------+
```

To update the status of a `todo` (e.g. change it to `completed` status), use the `todo` ID

```bash
./todo --update 5e9fd6b1bcd2fa6bd267d4c4,completed
```

List only the completed `todo`s

```bash
./todo --list completed
```

You should see the one you just updated

```bash
+----------------------------+--------------------------------+-----------+
|             ID             |          DESCRIPTION           |  STATUS   |
+----------------------------+--------------------------------+-----------+
| "5e9fd6b1bcd2fa6bd267d4c4" | Create an Azure Cosmos DB      | completed |
|                            | database account               |           |
+----------------------------+--------------------------------+-----------+
```

## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
