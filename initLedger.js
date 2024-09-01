const axios = require("axios");
const _ = require("lodash");

async function initializeLedger() {
  const apiUrl = "http://localhost:80/api/invoke/createAsset";
  const totalAssets = 1000;
  const batchSize = 500;
  const delayBetweenBatches = 1000;

  try {
    for (let i = 0; i < totalAssets; i += batchSize) {
      const payloads = generatePayloads(batchSize, i);

      const response = await sendBatchRequest(apiUrl, payloads);

      console.log(`Batch ${i + 1}-${i + payloads.length} sent successfully.`);

      if (i + batchSize < totalAssets) {
        await delay(delayBetweenBatches);
      }
    }

    console.log("Initialization completed successfully.");
  } catch (error) {
    console.error("Error initializing ledger:", error.message);
  }
}

function generatePayloads(batchSize, startIndex) {
  const payloads = [];

  for (let j = 0; j < batchSize; j++) {
    const projectName = `Powerplant ${startIndex + j + 1}`; // Generate project name based on index
    const amount = _.random(50, 200);
    const claimAmount = _.random(100, 300);
    const lastAction = `Initialized ${projectName} at ${getCurrentTime()}`;

    const payload = {
      "@assetType": "project",
      project: projectName,
      amount: amount,
      claimAmount: claimAmount,
      lastAction: lastAction,
    };

    payloads.push(payload);
  }

  return payloads;
}

async function sendBatchRequest(url, payloads) {
  try {
    const response = await axios.post(
      url,
      { asset: payloads },
      {
        headers: {
          "Content-Type": "application/json",
          "cache-control": "no-cache",
        },
      }
    );

    return response;
  } catch (error) {
    if (error.response) {
      console.error(
        `Request failed with status ${error.response.status}:`,
        error.response.data
      );
    } else {
      console.error("Request failed:", error.message);
    }
    throw error;
  }
}

function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function getCurrentTime() {
  const date = new Date();
  const localtime = new Date(date.getTime() + 8 * 60 * 60000); // Adjust for Malaysia timezone (UTC+8)
  return localtime.toLocaleString("en-US", { timeZone: "Asia/Kuala_Lumpur" });
}

initializeLedger();
