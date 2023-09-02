# dataloadgen-example
An example of how to use https://github.com/vikstrous/dataloadgen with gqlgen



This example was created by 

```
.
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── gqlgen.yml
├── graph
│   ├── generated.go
│   ├── loader
│   │   └── loader.go
│   ├── model
│   │   ├── models_gen.go
│   │   └── todo.go
│   ├── resolver.go
│   ├── schema.graphqls
│   ├── schema.resolvers.go
│   └── storage
│       └── storage.go
├── server.go
└── tools.go
```

```graphql
mutation{
  t1: createTodo(input:{text:"todo1",userId:"alice"}){
    id
    text
    done
    user{
      id
      name
    }
  }
  t2: createTodo(input:{text:"todo2",userId:"alice"}){
    id
    text
    done
    user{
      id
      name
    }
  }
  t3: createTodo(input:{text:"todo3",userId:"bob"}){
    id
    text
    done
    user{
      id
      name
    }
  }
}
```
```json
{
  "data": {
    "t1": {
      "id": "9014147064985197323",
      "text": "todo1",
      "done": false,
      "user": {
        "id": "alice",
        "name": "Alice"
      }
    },
    "t2": {
      "id": "763913058984159819",
      "text": "todo2",
      "done": false,
      "user": {
        "id": "alice",
        "name": "Alice"
      }
    },
    "t3": {
      "id": "5828640345075959780",
      "text": "todo3",
      "done": false,
      "user": {
        "id": "bob",
        "name": "Bob"
      }
    }
  }
}
```

```graphql
{
  todos{
    id
    text
    done
    user{
      id
      name
    }
  }
}
```
```json
{
  "data": {
    "todos": [
      {
        "id": "9014147064985197323",
        "text": "todo1",
        "done": false,
        "user": {
          "id": "alice",
          "name": "Alice"
        }
      },
      {
        "id": "763913058984159819",
        "text": "todo2",
        "done": false,
        "user": {
          "id": "alice",
          "name": "Alice"
        }
      },
      {
        "id": "5828640345075959780",
        "text": "todo3",
        "done": false,
        "user": {
          "id": "bob",
          "name": "Bob"
        }
      }
    ]
  }
}
```