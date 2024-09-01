const axios = require("axios");

// Create a shared Axios instance
const apiClient = axios.create({
  baseURL: "http://localhost/api",
  headers: {
    "Content-Type": "application/json",
    accept: "*/*",
  },
});

async function searchProjects() {
  try {
    const requestData = {
      query: {
        selector: {
          "@assetType": "project",
        },
      },
    };
    const response = await apiClient.post("/query/search", requestData);
    return response.data.result.map((project) => project["@key"]);
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
}

async function deleteLedger() {
  try {
    const keysArray = await searchProjects();

    // Construct payloads for each key
    const payloads = keysArray.map((key) => ({
      key: {
        "@assetType": "project",
        "@key": key,
      },
    }));

    // Make DELETE requests for each payload, controlling concurrency
    await Promise.all(
      payloads.map(async (payload) => {
        try {
          const response = await apiClient.delete("/invoke/deleteAsset", {
            data: payload,
          });
          console.log("API response:", response.data);
        } catch (error) {
          console.error("Error deleting asset:", error.message);
        }
      })
    );
  } catch (error) {
    console.error("Error searching projects:", error.message);
  }
}

deleteLedger();
