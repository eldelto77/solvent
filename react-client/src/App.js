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

  removeItem = item => {
    this.state.toDoList.removeItem(item.id);
    return this.setState({ toDoList: this.state.toDoList });
  }

  moveItem = (id, targetIndex) => {
    this.state.toDoList.moveItem(id, targetIndex);
    return this.setState({ toDoList: this.state.toDoList });
  }

  renameItem = (item, title) => {
    this.state.toDoList.renameItem(item.id, title);
    return this.setState({ toDoList: this.state.toDoList });
  }

  render() {
    return (
      <div className="App">
        <header>
          <div className="Header">
            <span className="HeaderSpacer"></span>
            <h1 className="HeaderTitle">Solvent</h1>
            <button className="HeaderMenu">_</button>
          </div>
        </header>
        <RToDoList
          toDoList={this.state.toDoList}
          checkItem={this.checkItem}
          addItem={this.addItem}
          removeItem={this.removeItem}
          moveItem={this.moveItem}
          renameItem={this.renameItem}
        />
      </div>
    );
  }
}

export default App;
