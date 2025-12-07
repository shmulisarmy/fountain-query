import { createSignal, type Component } from 'solid-js';
import { createMutable } from 'solid-js/store';



type RemoteUpdate = {
    Type: "add" | "remove" | "update"
    Data: any
    Path: string
    Source_name: string
}


type Person={
  name:string
  email:string
  id:number
  todo:{[key: string]: {
      epic_title:string
      author:string
      id:number
    }}
}



type Todo = Person['todo'][0]
function syncMessagesInto(Source_name: string, receiver: {}){
  return function(event: MessageEvent){

    // console.log(event.data)
    const update: RemoteUpdate = JSON.parse(event.data)
    console.log(update)
    switch (update.Type) {
      case "add":
        console.log("add", update.Data)
        let current = receiver
        const all_expect_first_and_last = update.Path.split("/").slice(1, -1)
        const last_key = update.Path.split("/").pop()
        console.log({all_expect_first_and_last})
        for (const key of all_expect_first_and_last) {
          current = current[key as keyof typeof current]
        } 
        console.log({current})
        console.log({data: update.Data})
        console.log({last_key})
        current[last_key as keyof typeof current] = JSON.parse(update.Data) as never
        break
        case "remove":
          console.log("remove", update.Data)
          break
          case "update":
            console.log("update", update.Data)
            break
            default:
              console.log("unknown type", update)
              console.log("unknown type")
              break
            }
          }
}

let source_name_upto = 0
function generate_source_name(){
  return "source_" + source_name_upto++
}
function listen_on(streaming_source: string, source_name: string = generate_source_name()): {} {
  const treeShapedReceiver = createMutable({})
  const ws = new WebSocket(streaming_source)
  ws.onmessage = syncMessagesInto(source_name, treeShapedReceiver)
  return treeShapedReceiver
}

const people: {[key: string]: Person} = listen_on("ws://localhost:8080/stream-data")


{
  (window as any).people = people
}



export function Todo_c({props}: {props:Todo}){
  return (
    <li class='flex flex-col space-y-2 ml-2 '>
      <p class="text-lg font-bold">{props.epic_title}</p>
      <p class="text-gray-500">{props.author}</p>
      <p class="text-xs text-gray-600">ID: {props.id}</p>
    </li>
  )
}
export function Person_c({props}: {props:Person}){
  return (
    <div class='flex flex-col space-y-2 p-1 m-1 border min-h-24'>
      <p class="text-lg font-bold">{props.name}</p>
      <p class="text-gray-500">{props.email}</p>
      <p class="text-xs text-gray-600">ID: {props.id}</p>
      <ul class='ml-2'>
        {Object.entries(props.todo).map(([id, todo]) => (
          <Todo_c  props={todo}/>
        ))}
      </ul>
    </div>
  )
}

const App: Component = () => {
  return (
    <div>
      {JSON.stringify(people)}
      <ul>
        {Object.entries(people).map(([id, person]) => (
          <Person_c  props={person}/>
        ))}
      </ul>
    </div>
  );
};

export default App;
