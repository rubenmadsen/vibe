# Costner

## Idea
"Costner" is an idea of an API testing application that, instead of the traditional API testing apps, like Postman and Insomnia, uses a graph based interface linking different nodes togeather. There should be different nodes available and the user can link them togeather by dragging one node's outputs to another node's input.

## Nodes
There should be different types of nodes. Here is a description of a few with purpose, input, output and gui elements
- envNode 
    1. Description - loads environment variables from the operating system or from .env files
    2. Interface - checkboxes for what variables should be imported, operating system environment variables and .env
    3. Output - all variables found
- requestNode
    1. Description - calls an endpoint
    2. Interface - url, headers, etc
    3. Output - the parts of the response that are relevand, headers, params, path, body etc..

## Example of use by a user
1. User adds an envNode to the canvas
2. User checks the operating variables checkbox (OS)
3. User add a requestNode to the canvas
4. user edits url
5. User drags from envNode to requestNode to make a specific variable available
6. User clicks "run" on the request causing its "dependancies to be called first(the envNode)

## technologies
- Golang
- Fyne

## Purpose
A way to visualize api testing flows

## Additions
- Project structure: Modular structure with separate packages for different components (UI, nodes, graph logic)
- Additional node types: transformNode (for data manipulation), conditionalNode (for branching logic), variableNode (for defining where variables should go in requests)
- Data flow: Structured Go types with JSON serialization for inter-node communication
- UI features: Basic drag-and-drop functionality (advanced features like zoom/pan for future iterations)
- Persistence: Nodes and graph configurations stored as JSON files for CLI compatibility - allows running without GUI
- CLI support: Application should be runnable from command line without GUI for automation

