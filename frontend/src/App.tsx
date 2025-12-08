import { createSignal, type Component } from 'solid-js';
import { live_db } from './live_db';






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


const people: {[key: string]: Person} = live_db("ws://localhost:8080/stream-data")





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
