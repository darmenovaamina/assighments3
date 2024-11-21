import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { getItems, createItem, deleteItem } from '../api'; // Ensure deleteItem is imported from the API

// Async thunk actions for fetching, adding, and deleting items
export const fetchItems = createAsyncThunk('items/fetchItems', async () => {
    const response = await getItems();
    return response.data;
});

export const addItem = createAsyncThunk('items/addItem', async (newItem) => {
    await createItem(newItem);
    const response = await getItems(); // Refresh list after adding item
    return response.data;
});

// New removeItem thunk for deleting an item
export const removeItem = createAsyncThunk('items/removeItem', async (id) => {
    await deleteItem(id);
    const response = await getItems(); // Refresh list after deleting item
    return response.data;
});

// Define the items slice
const itemsSlice = createSlice({
    name: 'items',
    initialState: {
        items: [],
        loading: false,
        error: null,
    },
    reducers: {},
    extraReducers: (builder) => {
        builder
            .addCase(fetchItems.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(fetchItems.fulfilled, (state, action) => {
                state.loading = false;
                state.items = action.payload;
            })
            .addCase(fetchItems.rejected, (state, action) => {
                state.loading = false;
                state.error = action.error.message;
            })
            .addCase(addItem.fulfilled, (state, action) => {
                state.items = action.payload;
            })
            .addCase(removeItem.fulfilled, (state, action) => {
                state.items = action.payload;
            });
    },
});

export default itemsSlice.reducer;
