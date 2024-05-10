# GoBank API

The GoBank API is a lightweight and efficient JSON API server implemented in Go. It provides endpoints for managing bank accounts and performing account-related operations. With easy-to-use RESTful endpoints, developers can quickly integrate banking functionalities into their applications.

## Features

- Handle GET requests to retrieve account information.
- Handle POST requests to create new accounts.
- Handle DELETE requests to delete existing accounts.
- Handle error responses gracefully.

## Setup

### 1. Clone the repository:

```bash
git clone https://github.com/nickemma/gobank-api.git
```

### 2. Navigate to the project directory:

```bash
 cd gobank-api
```

### 3. Build the project:

```bash
 go build / make build
```

### 4. Run the server:

```bash
 go run ./ OR make run
```

By default, the server listens on localhost:5000. You can customize the listen address by specifying the LISTEN_ADDR environment variable.

# API Endpoints

## Get Account

```bash
 GET /account/{id}
```

Retrieves account information based on the provided ID.

## Create Account

```bash
 POST /account
```

Creates a new account. Requires a JSON payload containing account details (e.g., first name, last name).

## Delete Account

```bash
 DELETE /account/{id}
```

Deletes an existing account based on the provided ID.

## Error Handling

The API server handles errors gracefully and returns appropriate error responses with descriptive error messages in JSON format.

## Dependencies

- Gorilla Mux: https://github.com/gorilla/mux

## 👤 Author <a name="author"></a>

👤 **Nicholas Emmanuel**

- GitHub: [@NickEmma](https://github.com/NickEmma)
- Twitter: [@techieEmma](https://twitter.com/techieEmma)
- LinkedIn: [@Nicholas Emmanuel](https://www.linkedin.com/in/techieemma/)

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## License

Please be advised that our project is released under the terms of a License. Please ensure that you read and understand the terms of the [MIT](LICENSE) License before using our project.

## Contact

### Nicholas Emmanuel

 <div align="center">
 <a href="https://www.linkedin.com/in/techieemma/"><img src="https://img.shields.io/badge/linkedin-%23f78a38.svg?style=for-the-badge&logo=linkedin&logoColor=white" alt="Linkedin"></a> 
 <a href="https://twitter.com/techieEmma"><img src="https://img.shields.io/badge/Twitter-%23f78a38.svg?style=for-the-badge&logo=Twitter&logoColor=white" alt="Twitter"></a> 
 <a href="https://github.com/nickemma/"><img src="https://img.shields.io/badge/github-%23f78a38.svg?style=for-the-badge&logo=github&logoColor=white" alt="Github"></a> 
 <a href="https://medium.com/@nicholasemmanuel321"><img src="https://img.shields.io/badge/Medium-%23f78a38.svg?style=for-the-badge&logo=Medium&logoColor=white" alt="Medium"></a> 
 <a href="mailto:nicholasemmanuel321@gmail.com"><img src="https://img.shields.io/badge/Gmail-f78a38?style=for-the-badge&logo=gmail&logoColor=white" alt="Linkedin"></a>
 </div>

## Acknowledgments

- [Creator](https://github.com/anthdm) for the inspiration for this project.
- [LazyCoders](https://lazy-coders.netlify.app/) For help and support throughout my development journey
