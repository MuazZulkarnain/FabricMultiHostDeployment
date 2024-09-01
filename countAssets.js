const axios = require("axios");

async function countAssets() {
  const url = "http://localhost/api/query/search";
  const requestData = {
    query: {
      selector: {
        "@assetType": "project",
      },
    },
  };

  try {
    const response = await axios.post(url, requestData);
    const assetCount = response.data.result.length; //count of assets
    console.log("No of assets in the ledger:", assetCount);
  } catch (error) {
    throw error.response ? error.response.data : error.message;
  }
}

countAssets();
