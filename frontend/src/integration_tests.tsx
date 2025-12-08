import { Person } from './App';
import { live_db } from './live_db';
import { deepObjectCompare } from './object_compare';
const backend_base_url = "localhost:8080"

export async function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

export function run_tests() {
  //integration tests
  // The whole idea of this project is that the data is the same, regardless of the order
  // in which things were updated or data was received in. Therefore, we check that here
  // by having two different clients that are supposed to represent the same data. 
  // although the first source is created (and starts receiving data before updates are triggered, it should have be the same as the second one
  const client_data_reflection_1: { [key: string]: Person; } = live_db(`ws://${backend_base_url}/stream-data`);
  fetch(`http://${backend_base_url}/add-person?name=donkey`).then(async function () {
    const client_data_reflection_2: { [key: string]: Person; } = live_db(`ws://${backend_base_url}/stream-data`);
    await sleep(200);
    const url = `https://json-diff-pro-copy-56b52272.base44.app/?actual=${encodeURIComponent(JSON.stringify(client_data_reflection_1))}&expected=${encodeURIComponent(JSON.stringify(client_data_reflection_2))}`;
    // `https://compare-production-1494.up.railway.app/compare?expected=${encodeURIComponent(JSON.stringify(client_data_reflection_1))}&actual=${encodeURIComponent(JSON.stringify(client_data_reflection_2))}`
    console.log({ url });
    if (!deepObjectCompare(client_data_reflection_1, client_data_reflection_2)) {
      alert("test failed");
      console.log(`to see what is different between the 2 json objects follow ${url}`);
      throw new Error("test failed");
    } else {
      alert("test passed");
    }
  });
}
