import React from 'react'
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd'

class RToDoList extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      newItemTitle: ""
    }
  }

  checkItem = item => {
    return this.props.checkItem(item);
  }

  addItem = event => {
    this.props.addItem(this.state.newItemTitle);
    this.setState({ newItemTitle: "" });
    event.preventDefault();
  }

  removeItem = item => {
    return this.props.removeItem(item);
  }

  setNewItemTitle = event => {
    this.setState({ newItemTitle: event.target.value });
  }

  moveItem = result => {
    const { destination, source, draggableId } = result;

    if (!destination) {
      return;
    }

    if (destination.droppableId === source.droppableId &&
      destination.index === source.index) {
      return;
    }

    this.props.moveItem(draggableId, destination.index);
  }

  renameItem = (item, event) => {
    return this.props.renameItem(item, event.target.value);
  }

  render() {
    return (
      <div className="ToDoList">
        <h1 className="ToDoListTitle">{this.props.toDoList.title}</h1>

        <div className="ToDoListBody">
          <RToDoItems
            items={this.props.toDoList.items}
            onCheck={this.checkItem}
            onRemove={this.removeItem}
            onDragEnd={this.moveItem}
            onRename={this.renameItem}
          />

          <AddItemBar
            value={this.state.newItemTitle}
            onChange={this.setNewItemTitle}
            onSubmit={this.addItem}
            disabled={this.state.newItemTitle.trim().length <= 0}
          />
        </div>
      </div>
    );
  }
}

function RToDoItems(props) {
  return (
    <DragDropContext onDragEnd={props.onDragEnd}>
      <Droppable droppableId="ToDoItemsDroppable">
        {provided => (
          <div className="ToDoItems" {...provided.droppableProps} ref={provided.innerRef}>
            {props.items.sort((a, b) => a.orderValue - b.orderValue)
              .map((item, index) => (
                <RToDoItem
                  key={item.id}
                  item={item}
                  index={index}
                  onCheck={props.onCheck}
                  onRemove={props.onRemove}
                  onRename={props.onRename}
                />
              ))}

            {provided.placeholder}
          </div>
        )}
      </Droppable>
    </DragDropContext>
  );
}

// Move to own component and keep newTitle in state.
// Apply renaming on focus loss
function RToDoItem(props) {
  return (
    <Draggable draggableId={props.item.id} index={props.index}>
      {provided => (
        <div
          className="ToDoItem"
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
        >
          <input
            className="ToDoItemCheckbox"
            type="checkbox"
            checked={props.item.checked}
            onChange={() => props.onCheck(props.item)}
          />
          <input
            className="ToDoItemTitle"
            type="text"
            value={props.item.title}
            onChange={event => props.onRename(props.item, event)}
          />
          <button
            className="ToDoItemDelete"
            onClick={() => props.onRemove(props.item)}>x</button>
        </div>
      )}
    </Draggable>
  );
}

function AddItemBar(props) {
  return (
    <form className="AddItemBar" onSubmit={props.onSubmit}>
      <span className="AddItemBarLogo">+</span>
      <input
        className="AddItemBarTitle"
        type="text"
        value={props.value}
        placeholder="New item"
        onChange={props.onChange}
      />
      <input className="AddItemBarButton" type="submit" value="Add" disabled={props.disabled} />
    </form>
  );
}

export default RToDoList;