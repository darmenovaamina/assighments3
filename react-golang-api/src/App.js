import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { fetchItems, addItem, removeItem } from "./store/itemsSlice";
import "./App.css";

function App() {
  const dispatch = useDispatch();
  const { items, loading, error } = useSelector((state) => state.items);
  const [newItem, setNewItem] = useState({ name: "", price: "" });

  // Fetch items on component mount
  useEffect(() => {
    dispatch(fetchItems());
  }, [dispatch]);

  const handleCreateItem = (e) => {
    e.preventDefault();
    // Dispatch addItem with the new item details
    dispatch(addItem({ name: newItem.name, price: parseFloat(newItem.price) }));
    setNewItem({ name: "", price: "" }); // Reset form fields
  };

  return (
    <div className="App">
      <h1>Items</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {loading ? (
        <p className="loading">Loading...</p>
      ) : (
        <ul>
          {items.map((item) => (
            <li key={item.id} style={{ display: "flex", justifyContent: "space-between" }}>
              <div style={{ alignContent: "center" }}>{item.name} - ${item.price}</div>
              <button
                style={{ margin: "0", background: "red" }}
                onClick={() => dispatch(removeItem(item.id))}
              >
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
      <form onSubmit={handleCreateItem}>
        <h2>Add New Item</h2>
        <input
          type="text"
          placeholder="Name"
          value={newItem.name}
          onChange={(e) => setNewItem({ ...newItem, name: e.target.value })}
          required
        />
        <input
          type="number"
          placeholder="Price"
          value={newItem.price}
          onChange={(e) => setNewItem({ ...newItem, price: e.target.value })}
          required
        />
        <button type="submit">Add Item</button>
      </form>
    </div>
  );
}

export default App;
