import React from 'react';
import './App.css';
import RToDoList from './solvent/render/RToDoList'
import ToDoList from './solvent/ToDoList'

class App extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      toDoList: ToDoList.new("List0")
    }

    this.state.toDoList.addItem("Item0")
  }

  checkItem = item => {
    if (item.checked) {
      this.state.toDoList.uncheckItem(item.id);
    } else {
      this.state.toDoList.checkItem(item.id);
    }
    return this.setState({ toDoList: this.state.toDoList });
  }

  addItem = title => {
    this.state.toDoList.addItem(title);
    return this.setState({ toDoList: this.state.toDoList });
  }

  render() {
    return (
      <div className="App">
        <RToDoList toDoList={this.state.toDoList} checkItem={this.checkItem} addItem={this.addItem} />
      </div>
    );
  }
}

export default App;
