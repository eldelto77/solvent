import React from 'react'

import AddItemBar from './AddItemBar'
import RToDoItems from './RToDoItems'

export default class RToDoList extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      newItemTitle: ""
    }
  }

  addItem = event => {
    this.props.addItem(this.state.newItemTitle);
    this.setState({ newItemTitle: "" });
    event.preventDefault();
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

  render() {
    return (
      <div className="ToDoList">
        <h1 className="ToDoListTitle">{this.props.toDoList.title}</h1>

        <div className="ToDoListBody">
          <RToDoItems
            items={this.props.toDoList.items}
            onCheck={this.props.checkItem}
            onRemove={this.props.removeItem}
            onDragEnd={this.moveItem}
            onRename={this.props.renameItem}
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
