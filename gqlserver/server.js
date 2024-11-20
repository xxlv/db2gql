const { ApolloServer, gql } = require("apollo-server");
const axios = require("axios");

// Retry configuration
const MAX_RETRIES = 10; // Maximum number of retries
const RETRY_DELAY_MS = 5000; // 5 seconds delay between retries

async function loadSchema() {
  let retries = 0;

  while (retries < MAX_RETRIES) {
    try {
      const response = await axios.get("http://localhost:8080/previewSchema");

      // Check if the response contains valid data
      if (response.data) {
        return gql(response.data);
      }

      console.log("No data received, retrying...");
    } catch (err) {
      console.error("Error fetching schema:", err.message);
    }

    retries++;
    console.log(`Retrying... attempt ${retries}/${MAX_RETRIES}`);
    await new Promise((resolve) => setTimeout(resolve, RETRY_DELAY_MS)); // Wait before retrying
  }

  throw new Error("Failed to load schema after multiple retries");
}

async function startServer() {
  try {
    const typeDefs = await loadSchema();
    const resolvers = {};

    const server = new ApolloServer({ typeDefs, resolvers });

    server.listen().then(({ url }) => {
      console.log(`Server ready at ${url}`);
    });
  } catch (err) {
    console.error("Error starting server:", err);
  }
}

startServer();
