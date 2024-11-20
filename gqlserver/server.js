const { ApolloServer, gql } = require("apollo-server");
const axios = require("axios");

async function loadSchema() {
  const response = await axios.get("http://localhost:8080/previewSchema");
  return gql(response.data);
}

async function startServer() {
  const typeDefs = await loadSchema();
  const resolvers = {};

  const server = new ApolloServer({ typeDefs, resolvers });

  server.listen().then(({ url }) => {
    console.log(`Server ready at ${url}`);
  });
}

startServer().catch((err) => {
  console.error("Error starting server:", err);
});
