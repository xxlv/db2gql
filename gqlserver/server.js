const { ApolloServer } = require("apollo-server");
const { GraphQLObjectType, buildSchema } = require("graphql");
const axios = require("axios");
const { faker } = require("@faker-js/faker");
const crypto = require("crypto");

const CHECK_INTERVAL_MS = 10000 / 2; // 5s
let schemaHash = "";
let server;

// Function to fetch the schema from the endpoint
async function fetchSchema() {
  try {
    const response = await axios.get("http://localhost:8080/previewSchema");
    if (!response.data || typeof response.data !== "string") {
      throw new Error("Schema data is empty or not a string");
    }
    return response.data;
  } catch (err) {
    console.error("Error fetching schema:", err.message);
    return null;
  }
}

// Advanced mock data generator
function generateMockData(type, depth = 0) {
  // Prevent infinite recursion
  if (depth > 2) return null;

  // Handle different type names
  switch (type.name) {
    case "String":
      return faker.lorem.word();
    case "Int":
      return faker.number.int();
    case "Float":
      return faker.number.float();
    case "Boolean":
      return faker.datatype.boolean();
    case "ID":
      return `mock-id-${faker.string.uuid()}`;
    default:
      // For complex types, create a nested mock object
      if (type.getFields) {
        const mockObj = {};
        const fields = type.getFields();

        Object.keys(fields).forEach((fieldName) => {
          let fieldType = fields[fieldName].type;

          // Unwrap non-null and list types
          while (fieldType.ofType) {
            fieldType = fieldType.ofType;
          }

          // Recursively generate mock data for nested fields
          mockObj[fieldName] = generateMockData(fieldType, depth + 1);
        });

        return mockObj;
      }

      return faker.lorem.sentence();
  }
}

// Function to create mock resolvers for all fields in the schema
function createMockResolvers(schema) {
  const mockResolvers = {};
  const typeMap = schema.getTypeMap();

  Object.keys(typeMap).forEach((typeName) => {
    const type = typeMap[typeName];

    // Skip introspection types
    if (typeName.startsWith("__")) return;

    // Check if type has fields
    if (type instanceof GraphQLObjectType && type.getFields) {
      mockResolvers[typeName] = {};
      const fields = type.getFields();

      Object.keys(fields).forEach((fieldName) => {
        mockResolvers[typeName][fieldName] = (parent, args, context, info) => {
          let fieldType = fields[fieldName].type;

          // Unwrap non-null and list types
          while (fieldType.ofType) {
            fieldType = fieldType.ofType;
          }

          // Generate mock data for the specific field
          return generateMockData(fieldType);
        };
      });
    }
  });

  return mockResolvers;
}

// Function to start Apollo server
async function startApolloServer() {
  try {
    // Fetch and build the schema
    const schemaString = await fetchSchema();
    if (!schemaString) {
      throw new Error("Failed to fetch schema");
    }

    if (schemaString.length <= 16) {
      console.log("schema not prepared...");
      return;
    }

    // Build the schema using graphql's buildSchema
    const schema = buildSchema(schemaString);

    // Create mock resolvers
    const mockResolvers = createMockResolvers(schema);

    // Start the server with the built schema
    server = new ApolloServer({
      schema,
      mocks: mockResolvers,
      mockEntireSchema: false,
    });

    const { url } = await server.listen();
    console.log(`Server ready at ${url}`);
  } catch (err) {
    console.error("Error starting server:", err.message);
    console.error("Full error:", err);
  }
}

// Function to restart the server
async function restartServer() {
  if (server) {
    await server.stop();
    console.log("Server stopped. Restarting...");
  }

  await startApolloServer();
}

// Function to monitor schema changes
async function monitorSchemaChanges() {
  setInterval(async () => {
    try {
      const schemaData = await fetchSchema();

      if (schemaData) {
        const newHash = crypto
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

// Initialize server
(async () => {
  await startApolloServer();
  monitorSchemaChanges();
})();
