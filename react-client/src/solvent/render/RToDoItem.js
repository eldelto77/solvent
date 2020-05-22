import React from 'react'
import { Draggable } from 'react-beautiful-dnd'

// Move to own component and keep newTitle in state.
// Apply renaming on focus loss
export default class RToDoItem extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      editing: false,
      newTitle: ""
    }
  }

  startEditing = () => {
    this.setState({
      editing: true,
      newTitle: this.props.item.title
    });
  }

  editItem = event => {
    this.setState({newTitle: event.target.value});
  }

  renameItem = () => {
    this.props.onRename(this.props.item, this.state.newTitle);
    this.setState({
      editing: false,
      newTitle: ""
    });
  }

  render() {
    return (
      <Draggable draggableId={this.props.item.id} index={this.props.index}>
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
              checked={this.props.item.checked}
              onChange={() => this.props.onCheck(this.props.item)}
            />
            <input
              className="ToDoItemTitle"
              type="text"
              value={this.state.editing ? this.state.newTitle : this.props.item.title}
              onFocus={this.startEditing}
              onChange={this.editItem}
              onBlur={this.renameItem}
            />
            <button
              className="ToDoItemDelete"
              onClick={() => this.props.onRemove(this.props.item)}>x</button>
          </div>
        )}
      </Draggable>
    );
  }
}