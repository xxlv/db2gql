const { ApolloServer, gql } = require("apollo-server");
const axios = require("axios");

const CHECK_INTERVAL_MS = 10000; // 10s
let schemaHash = "";
let server;

async function fetchSchema() {
  try {
    const response = await axios.get("http://localhost:8080/previewSchema");
    return response.data || null;
  } catch (err) {
    console.error("Error fetching schema:", err.message);
    return null;
  }
}

async function loadSchema() {
  const schemaData = await fetchSchema();

  if (schemaData) {
    return gql(schemaData);
  }

  throw new Error("Failed to load schema");
}

async function startApolloServer() {
  try {
    const typeDefs = await loadSchema();
    const resolvers = {};

    server = new ApolloServer({ typeDefs, resolvers });

    const { url } = await server.listen();
    console.log(`Server ready at ${url}`);
  } catch (err) {
    console.error("Error starting server:", err);
  }
}

async function restartServer() {
  if (server) {
    await server.stop();
    console.log("Server stopped. Restarting...");
  }

  await startApolloServer();
}

async function monitorSchemaChanges() {
  setInterval(async () => {
    try {
      const schemaData = await fetchSchema();

      if (schemaData) {
        const newHash = require("crypto")
          .createHash("md5")
          .update(schemaData)
          .digest("hex");

        if (newHash !== schemaHash) {
          console.log("Schema has changed. Restarting server...");
          schemaHash = newHash;
          await restartServer();
        }
      }
    } catch (err) {
      console.error("Error monitoring schema changes:", err.message);
    }
  }, CHECK_INTERVAL_MS);
}

(async () => {
  await startApolloServer();
  monitorSchemaChanges();
})();
