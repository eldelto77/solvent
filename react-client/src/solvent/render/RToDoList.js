import React from 'react'

import TextareaAutosize from 'react-textarea-autosize';

import AddItemBar from './AddItemBar'
import RToDoItems from './RToDoItems'

export default class RToDoList extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      editing: false,
      newTitle: this.props.toDoList.title.value,
      newItemTitle: ""
    }
  }

  startTitleEditing = () => {
    this.setState({
      editing: true,
      newTitle: this.props.toDoList.title.value
    });
  }

  editListTitle = event => {
    this.setState({ newTitle: event.target.value });
  }

  renameList = () => {
    this.props.renameList(this.state.newTitle);
    this.setState({
      editing: false,
      newTitle: this.props.toDoList.title.value
    });
  }

  addItem = event => {
    this.props.addItem(this.state.newItemTitle);
    this.setState({ newItemTitle: "" });
    event.preventDefault();
  }

  editNewItemTitle = event => {
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
      <div className={"ToDoList" + (this.props.toDoList.isChecked() ? " checked" : "")}>
        <TextareaAutosize
          className="ToDoListTitle"
          placeholder="Title"
          onFocus={this.startTitleEditing}
          onChange={this.editListTitle}
          onBlur={this.renameList}
        >
          {this.state.editing ? this.state.newTitle : this.props.toDoList.title.value}
        </TextareaAutosize>

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
            onChange={this.editNewItemTitle}
            onSubmit={this.addItem}
            disabled={this.state.newItemTitle.trim().length <= 0}
          />
        </div>
      </div>
    );
  }
}
