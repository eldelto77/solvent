import React from 'react'

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

  setNewItemTitle = event => {
    this.setState({ newItemTitle: event.target.value });
  }

  addItem = event => {
    this.props.addItem(this.state.newItemTitle);
    this.setState({ newItemTitle: "" });
    event.preventDefault();
  }

  render() {
    return (
      <div className="ToDoList">
        <h1>{this.props.toDoList.title}</h1>

        <RToDoItems items={this.props.toDoList.items} onCheck={this.checkItem} />

        <AddItemBar value={this.state.newItemTitle} onChange={this.setNewItemTitle} onSubmit={this.addItem} />
      </div>
    );
  }
}

function RToDoItems(props) {
  return (
    props.items.sort((a, b) => a.orderValue - b.orderValue).map(item => {
      return (
        <div key={item.id}>
          <input type="checkbox" checked={item.checked} onChange={() => props.onCheck(item)} />
          {item.title}
        </div>
      );
    })
  );
}

function AddItemBar(props) {
  return (
    <form onSubmit={props.onSubmit}>
      <input type="text" value={props.value} onChange={props.onChange} />
      <input type="submit" value="Add" />
    </form>
  );
}

export default RToDoList;