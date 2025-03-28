Example of a Simple Web Application in Go

This is an example of a simple web application written in the Go programming language. The application is an HTTP server that handles requests and renders HTML templates.

Installation

Clone the repository:

git clone https://gitea.kood.tech/artjomkulikovski/art

Navigate to the project directory:

cd your/repository

Run the application:

go run main.go

Usage

After starting the application, open a web browser and navigate to http://localhost:8080.
You will see a web page with a form for entering data.
Enter your data and select the processing method (decode or encode).
Submit the form to see the processing result.

Project Structure

main.go: The main application file containing the logic of the HTTP server.
template/: Directory containing HTML templates.
static/: Directory containing static files (CSS).
Route Handlers
/: Handler for the root route, which displays the main page of the application.
/decode: Handler for processing user-entered data with subsequent decoding or encoding.

Functions

renderTemplate: This function takes an http.ResponseWriter, a string name representing the name of the template to render, and a Data struct containing the data to be passed to the template. It retrieves the template from the templates map by name and renders it to the response writer. If the template does not exist, it returns an internal server error.

decodeInput: This function takes a string input representing the user's input data. It checks for errors in the input data using the checkErrors function. If there are no errors, it splits the input into lines and processes each line using the processLine function. It returns a slice of strings representing the decoded input and a status code indicating success.

encodeInput: This function takes a string input representing the user's input data. It checks for errors in the input data using the checkErrors function. If there are no errors, it splits the input into lines and processes each line using the encodeLine function. It returns a slice of strings representing the encoded input and a status code indicating success.

checkErrors: This function takes a string input representing the user's input data. It uses regular expressions to check for various errors in the input data, such as invalid characters, invalid input format, invalid character count, missing characters, and unbalanced brackets. If any errors are found, it returns an error message. Otherwise, it returns nil indicating no errors.
processLine: This function takes a string line representing a single line of input data. It processes the line by replacing encoded segments with the corresponding characters. It returns the processed line as a string.

isBalanced: This function takes a string s representing a sequence of characters. It checks if the sequence has balanced brackets (i.e., every opening bracket has a corresponding closing bracket). It returns true if the brackets are balanced and false otherwise.

splitCountAndChars: This function takes a string s representing a count followed by characters. It splits the string into two parts: the count and the characters. If there is no space in the string, it returns an empty string for the count. It returns a slice of strings containing the count and characters.
startsWithNumber: This function takes a string s representing a sequence of characters. It checks if the sequence starts with a number. It returns true if the sequence starts with a number and false otherwise.