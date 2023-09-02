# dataloadgen-example
This is an example of how to use [vikstrous/dataloadgen](https://github.com/vikstrous/dataloadgen) with [gqlgen](https://gqlgen.com/) to improve graphql performance by batching and caching requests to the underlying storage system.

This example was created by following the [official tutorial](https://gqlgen.com/getting-started/). After completing the tutorial, an example in-memory storage package was created and a loader was created for the User object.

This file listing highlights the important files to look at to understand how loaders can be used with gqlgen.
```
.
├── gqlgen.yml
├── graph
│   ├── generated.go
│   ├── loader
│   │   └── loader.go - the implementation of loaders
│   ├── model
│   │   ├── models_gen.go
│   │   └── todo.go
│   ├── resolver.go
│   ├── schema.graphqls
│   ├── schema.resolvers.go - the implementation of resolvers that use loaders
│   └── storage
│       └── storage.go - the underlying storage system with artificial delays and logging
├── server.go - the wiring of the stoage system, loaders middleware and resolvers
└── tools.go
```

To use a loader, the recommended pattern is to create a new one for every HTTP request and inject it into the context using a middleware. [server.go](https://github.com/vikstrous/dataloadgen-example/blob/master/server.go) contains the wiring for this. Note the call to `loader.Middleware`. That allows resolvers to access the loader using `loader.Get(ctx)` and then call the methods on the loader objects.

## Try it out

Start the server by running:
```
go run .
```

Go to http://localhost:8080 and execute the following query to populate the TODOs:

```graphql
mutation {
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
The response should look like:
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
It intentionally takes several seconds to execute the mutations. Every access to the storage package is artificially delayed by a second and accesses are logged to stdout. The output in the console after running this query should look like:
```
UserStorage.Get
TodoStorage.Put
UserStorage.Get
TodoStorage.Put
UserStorage.Get
TodoStorage.Put
```

Note that there are only three calls to `UserStorage.Get` even though the user is fetched both in the execution of the mutation and later in the query for the user. This is because, in the root resolver, after accessing the user storage, the data loader cache is primed. See [schema.resolvers.go](https://github.com/vikstrous/dataloadgen-example/blob/master/graph/schema.resolvers.go) `mutationResolver.CreateTodo()` for how the cache is primed.

Now make another, read-only query to fetch all the todos along with their associated users.

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
The response should look like:
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
The output in the console should look like:
```
TodoStorage.GetAll
UserStorage.GetMulti 2
```

There is no call to `UserStorage.Get` in this case. The use of a loader in [schema.resolvers.go](https://github.com/vikstrous/dataloadgen-example/blob/master/graph/schema.resolvers.go) in `todoResolver.User()` causes concurrent executions to be batched, deduplicated and cached, so only a single call to `GetMulti` is made instead and only with 2 user IDs.
