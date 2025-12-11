import { createSignal, type Component, For, Show } from 'solid-js';
import { live_db as live_db_view } from "../../live_db_sdks/solid/live_db";

export type Person = {
  name: string
  email: string
  age: number
  id: number
  profile_picture?: string
}

export const backend_base_url = "localhost:8080"

const people: { [key: string]: Person } = live_db_view(`ws://${backend_base_url}/stream-data`);

// Documentation content
const DOCUMENTATION = {
  overview: {
    title: "How Real-Time Sync Works",
    content: `The Live DB system ensures that all clients always see the same data, regardless of when they connect or in what order updates occur.

**Key Principles:**
1. **WebSocket Streaming**: Data flows from server to client via WebSocket messages
2. **Event-Driven Updates**: Each change (add/remove/update) is broadcast as a sync message
3. **Path-Based Hierarchy**: Data structure is maintained through hierarchical paths
4. **Idempotent Operations**: Operations can be applied in any order with the same result

**Sync Message Types:**
- **Load**: Initial full data snapshot when connecting
- **Add**: New item added at a specific path
- **Remove**: Item deleted from a path
- **Update**: Existing item modified

The backend's EventEmitterTree automatically tracks all changes to SQL query results and broadcasts them to all connected clients in real-time.`
  },
  add: {
    title: "Add Operations",
    content: `When you add a new person, here's what happens:

1. **Frontend Request**: Your button click sends an HTTP request to the backend
2. **Backend Processing**: The server inserts the new row into the database table
3. **Query Re-evaluation**: The SQL query automatically re-runs (SELECT person...)
4. **Diff Detection**: The EventEmitterTree compares old vs new query results
5. **Sync Message**: An "add" message is sent via WebSocket to all connected clients
6. **Client Update**: Each client applies the update at the specified path

**Why It Stays in Sync:**
- All clients receive the same sync messages in order
- Path-based updates ensure data structure consistency
- The hierarchical structure matches your nested SQL query structure`
  },
  filter: {
    title: "Filtering",
    content: `The filtering shown here is client-side, but the underlying data sync is server-driven.

**Current Implementation:**
- Filters apply to the already-synced data in your browser
- All filtered results come from the same synchronized dataset
- Real-time updates continue to flow even when filtering

**Server-Side Filtering (SQL WHERE clause):**
The backend query includes: \`WHERE person.age >= 3\`

This means:
- Only persons with age >= 3 are included in the sync stream
- If a person's age changes below 3, they'll be removed from all clients
- If a new person with age >= 3 is added, they'll appear on all clients

**Client-Side Filtering (UI):**
- You can further filter the already-synced data
- This is for presentation only - doesn't affect the sync mechanism
- All clients have the same base data, but can display different filtered views`
  },
  delete: {
    title: "Delete Operations",
    content: `When you delete a person:

1. **HTTP Request**: GET request sent to backend delete endpoint
2. **Table Update**: Row removed from database table
3. **Query Re-evaluation**: SQL query automatically detects the change
4. **Remove Message**: "remove" sync message broadcast to all connected clients
5. **Path Removal**: Each client removes the person at their path (e.g., /person_123)

**Real-Time Sync:**
- All connected clients receive the same "remove" message simultaneously
- The person disappears from all browser windows at the same time
- Late-connecting clients will never see the deleted person (they get the current state)`
  },
  load: {
    title: "Load Sample Data",
    content: `Loading sample data triggers multiple operations:

1. **Bulk Insert**: Multiple persons inserted into tables
2. **Batch Processing**: All inserts processed together
3. **Query Updates**: SQL query detects all new matching rows
4. **Multiple Sync Messages**: Each new item generates an "add" message
5. **Client Reconciliation**: Clients apply all updates in order

**Consistency Guarantee:**
Even if you connect two clients at different times:
- Client A connects before data is loaded
- Client B connects after data is loaded
- Both will receive the same final state
- The Load message ensures late-connecting clients get full snapshot
- All subsequent sync messages apply identically`
  },
  };


function Person_c({ props, onDelete }: { props: Person; onDelete: (id: number) => void; }) {
  const profilePicture = () => props.profile_picture || `https://api.dicebear.com/7.x/avataaars/svg?seed=${props.name}`;

  return (
    <div class='flex flex-col space-y-3 p-4 m-2 border-2 border-gray-300 rounded-lg bg-white shadow-sm hover:shadow-md transition-shadow min-h-32'>
      <div class="flex justify-between items-start gap-3">
        <img
          src={profilePicture()}
          alt={`${props.name}'s profile`}
          class="w-16 h-16 rounded-full border-2 border-gray-300"
        />
        <div class="flex-1">
          <p class="text-xl font-bold text-gray-900">{props.name}</p>
          <p class="text-sm text-gray-600">{props.email}</p>
          <p class="text-xs text-gray-500">Age: {props.age} | ID: {props.id}</p>
        </div>
        <div class="flex gap-2">
          <button
            onClick={() => onDelete(props.id)}
            class="px-3 py-1 bg-red-500 text-white text-sm rounded hover:bg-red-600 transition-colors"
            title="Delete person"
          >
            Delete
          </button>
        </div>
      </div>
          </div>
  )
}

const App: Component = () => {
  const [filterName, setFilterName] = createSignal("");
  const [filterEmail, setFilterEmail] = createSignal("");
  const [activeDoc, setActiveDoc] = createSignal<keyof typeof DOCUMENTATION | null>(null);
  const [showAddPersonForm, setShowAddPersonForm] = createSignal(false);
    const [newPersonName, setNewPersonName] = createSignal("");
  const [newPersonEmail, setNewPersonEmail] = createSignal("");
    const [notificationMessage, setNotificationMessage] = createSignal<string | null>(null);

  // Get filtered people
  const filteredPeople = () => {
    const nameFilter = filterName().toLowerCase();
    const emailFilter = filterEmail().toLowerCase();
    return Object.entries(people).filter(([_, person]) => {
      const nameMatch = !nameFilter || person.name.toLowerCase().includes(nameFilter);
      const emailMatch = !emailFilter || person.email.toLowerCase().includes(emailFilter);
      return nameMatch && emailMatch;
    });
  };

  const handleAddPerson = async () => {
    const name = newPersonName().trim();
    const email = newPersonEmail().trim() || `${name.toLowerCase()}@example.com`;
    if (!name) return;
    
    await fetch(`http://${backend_base_url}/add-person?name=${encodeURIComponent(name)}&email=${encodeURIComponent(email)}`);
    setNewPersonName("");
    setNewPersonEmail("");
    setShowAddPersonForm(false);
    setActiveDoc("add");
  };

  
  
  

  
  const handleLoadSampleData = async () => {
    await fetch(`http://${backend_base_url}/add-sample-data`);
    setActiveDoc("load");
  };

  const handleDeletePerson = async (personId: number) => {
    if (!confirm(`Delete person with ID ${personId}?`)) {
      return;
    }
    await fetch(`http://${backend_base_url}/delete-person?id=${personId}`);

    setNotificationMessage(`🗑️ Person deleted! The frontend just made a request to delete the person. Because the data is subscribed through a query, the person automatically disappears in real-time across all connected clients. Open another tab to test the real-time sync!`);
    setTimeout(() => setNotificationMessage(null), 10000);
    setActiveDoc("delete");
  };

  return (
    <div class="min-h-screen bg-gray-100">
      {/* Header */}
      <header class="bg-blue-600 text-white shadow-lg">
        <div class="container mx-auto px-4 py-4">
          <h1 class="text-3xl font-bold">Live SQL Database Demo</h1>
          <p class="text-blue-100 mt-1">Real-time synchronized data across all clients</p>
        </div>
      </header>

      <div class="container mx-auto px-4 py-6">
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Content Area */}
          <div class="lg:col-span-2 space-y-4">
            {/* Notification Banner */}
            <Show when={notificationMessage()}>
              <div class="bg-blue-50 border-l-4 border-blue-500 p-4 rounded-lg shadow-md">
                <div class="flex justify-between items-start">
                  <p class="text-sm text-blue-900">{notificationMessage()}</p>
                  <button
                    onClick={() => setNotificationMessage(null)}
                    class="text-blue-700 hover:text-blue-900 ml-4"
                  >
                    ✕
                  </button>
                </div>
              </div>
            </Show>

            {/* Action Buttons */}
            <div class="bg-white rounded-lg shadow p-4">
              <h2 class="text-xl font-semibold mb-4 text-gray-800">Actions</h2>
              <div class="flex flex-wrap gap-3">
                <button
                  onClick={() => { setShowAddPersonForm(true); setActiveDoc("add"); }}
                  class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition-colors font-medium"
                >
                  ➕ Add Person
                </button>
                <button
                  onClick={handleLoadSampleData}
                  class="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600 transition-colors font-medium"
                >
                  🔄 Load Sample Data
                </button>
                <button
                  onClick={() => setActiveDoc("overview")}
                  class="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition-colors font-medium"
                >
                  📚 Show Sync Docs
                </button>
              </div>

              {/* Add Person Form */}
              <Show when={showAddPersonForm()}>
                <div class="mt-4 p-4 bg-green-50 rounded border border-green-200">
                  <h3 class="font-semibold mb-2">Add New Person</h3>
                  <div class="space-y-2">
                    <input
                      type="text"
                      placeholder="Name"
                      value={newPersonName()}
                      onInput={(e) => setNewPersonName(e.currentTarget.value)}
                      class="w-full px-3 py-2 border rounded"
                      onKeyDown={(e) => e.key === 'Enter' && handleAddPerson()}
                    />
                    <input
                      type="email"
                      placeholder="Email (optional)"
                      value={newPersonEmail()}
                      onInput={(e) => setNewPersonEmail(e.currentTarget.value)}
                      class="w-full px-3 py-2 border rounded"
                      onKeyDown={(e) => e.key === 'Enter' && handleAddPerson()}
                    />
                    <div class="flex gap-2">
                      <button
                        onClick={handleAddPerson}
                        class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700"
                      >
                        Add
                      </button>
                      <button
                        onClick={() => { setShowAddPersonForm(false); setNewPersonName(""); setNewPersonEmail(""); }}
                        class="px-4 py-2 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                </div>
              </Show>

                          </div>

            {/* Filters */}
            <div class="bg-white rounded-lg shadow p-4">
              <h2 class="text-xl font-semibold mb-4 text-gray-800">Filters</h2>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Filter by Name</label>
                  <input
                    type="text"
                    placeholder="Search names..."
                    value={filterName()}
                    onInput={(e) => { setFilterName(e.currentTarget.value); setActiveDoc("filter"); }}
                    class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Filter by Email</label>
                  <input
                    type="text"
                    placeholder="Search emails..."
                    value={filterEmail()}
                    onInput={(e) => { setFilterEmail(e.currentTarget.value); setActiveDoc("filter"); }}
                    class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            </div>

            {/* Data Display */}
            <div class="bg-white rounded-lg shadow p-4">
              <h2 class="text-xl font-semibold mb-4 text-gray-800">
                People ({filteredPeople().length} of {Object.keys(people).length})
              </h2>
              <Show
                when={filteredPeople().length > 0}
                fallback={<p class="text-gray-500 text-center py-8">No people match your filters</p>}
              >
                <div class="grid grid-cols-1 gap-4">
                  <For each={filteredPeople()}>
                    {([_, person]) => <Person_c props={person} onDelete={handleDeletePerson} />}
                  </For>
                </div>
              </Show>
            </div>
          </div>

          {/* Documentation Sidebar */}
          <div class="lg:col-span-1">
            <div class="bg-white rounded-lg shadow-lg p-4 sticky top-4">
              <div class="flex justify-between items-center mb-4">
                <h2 class="text-xl font-semibold text-gray-800">📖 Documentation</h2>
                <Show when={activeDoc()}>
                  <button
                    onClick={() => setActiveDoc(null)}
                    class="text-gray-500 hover:text-gray-700 text-sm"
                  >
                    ✕ Close
                  </button>
                </Show>
              </div>
              
              <Show
                when={activeDoc()}
                fallback={
                  <div class="text-gray-600 text-sm space-y-3">
                    <p class="font-medium">🚀 Experience Real-Time Sync!</p>
                    <div class="p-3 bg-yellow-50 border border-yellow-200 rounded">
                      <p class="font-semibold text-yellow-800">🧪 Try This Now:</p>
                      <p class="text-yellow-700 text-xs mt-1">Open this page in another tab or browser window, then make changes in one tab and watch them instantly appear in both!</p>
                    </div>
                    <ul class="list-disc list-inside space-y-2 ml-2">
                                            <li>Click delete buttons to see real-time removal across tabs</li>
                      <li>Use <strong>Filters</strong> to understand client-side vs server-side filtering</li>
                      <li>Click <strong>Load Sample Data</strong> to see batch updates</li>
                    </ul>
                    <button
                      onClick={() => setActiveDoc("overview")}
                      class="w-full mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
                    >
                      Show Full Documentation
                    </button>
                  </div>
                }
              >
                <div class="prose prose-sm max-w-none">
                  <h3 class="text-lg font-bold text-gray-900 mb-3">
                    {DOCUMENTATION[activeDoc()!]?.title}
                  </h3>
                  <div class="text-gray-700 whitespace-pre-line leading-relaxed">
                    {DOCUMENTATION[activeDoc()!]?.content}
                  </div>
                  <div class="mt-4 pt-4 border-t border-gray-200">
                    <p class="text-xs text-gray-500 mb-2">Related topics:</p>
                    <div class="flex flex-wrap gap-2">
                      <button
                        onClick={() => setActiveDoc("overview")}
                        class="text-xs px-2 py-1 bg-gray-100 rounded hover:bg-gray-200"
                      >
                        Overview
                      </button>
                                          </div>
                  </div>
                </div>
              </Show>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default App;
