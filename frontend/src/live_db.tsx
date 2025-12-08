import { createMutable,  } from 'solid-js/store';


function reactive_mutable_set(obj: any, receiver: {}) {
  for (const key of Object.keys(obj)) {
    receiver[key as keyof typeof receiver ] = obj[key] as keyof typeof receiver;
  }
}

type RemoteUpdate = {
  Type: "add" | "remove" | "update" | "load";
  Data: any;
  Path: string;
  Source_name: string;
};
function syncMessagesInto(receiver: {}) {
  return function (event: MessageEvent) {

    // console.log(event.data)
    const update: RemoteUpdate = JSON.parse(event.data);
    console.log(update);
    switch (update.Type) {
      case "add":
        console.log("add", update.Data);
        let current = receiver;
        const all_expect_first_and_last = update.Path.split("/").slice(1, -1);
        const last_key = update.Path.split("/").pop();
        console.log({ all_expect_first_and_last });
        for (const key of all_expect_first_and_last) {
          current = current[key as keyof typeof current];
        }
        console.log({ current });
        console.log({ data: update.Data });
        console.log({ last_key });
        current[last_key as keyof typeof current] = JSON.parse(update.Data) as never;
        break;
      case "remove":
        console.log("remove", update.Data);
        break;
      case "update":
        console.log("update", update.Data);
        break;
      case "load":
        reactive_mutable_set(JSON.parse(update.Data), receiver);
        break;
      default:
        console.log("unknown type", update);
        console.log("unknown type");
        break;
    }
  };
}
export function live_db(streaming_source: string): {} {
  const treeShapedReceiver = createMutable({});
  const ws = new WebSocket(streaming_source);
  ws.onmessage = syncMessagesInto(treeShapedReceiver);
  return treeShapedReceiver;
}
