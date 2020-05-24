import React from 'react';

import './App.css';
import DetailView from './solvent/render/DetailView'
import ListView from './solvent/render/ListView'

import ToDoList from './solvent/ToDoList'

import { toDoListFromDto, toDoListToDto } from './solvent/Dto'

class App extends React.Component {

  constructor(props) {
    super(props);

    const toDoList0 = ToDoList.new("List0");
    toDoList0.addItem("Item0");
    toDoList0.addItem("Item1");
    toDoList0.addItem("Item2");

    const toDoList1 = ToDoList.new("List1");
    toDoList1.addItem("Item3");
    toDoList1.addItem("Item4");
    toDoList1.addItem("Item5");

    this.state = {
      toDoLists: [toDoList0, toDoList1],
      activeToDoList: null,
      isListViewActive: true
    }
  }

  componentDidMount() {
    this.timer = setInterval(() => this.loadToDoLists(), 1000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
    this.timer = null;
  }

  loadToDoLists = () => {
    fetch("/api/to-do-list")
      .then(response => response.json())
      .then(responseBody => responseBody.toDoLists.map(toDoListFromDto))
      .then(toDoLists => {
        this.setState({ toDoLists: toDoLists });

        if (this.state.activeToDoList) {
          const activeToDoList = toDoLists.find(list => list.id === this.state.activeToDoList.id);
          if (activeToDoList) {
            this.setState({ activeToDoList: activeToDoList });
          }
        }
      });
    // Execute PUT calls on each change
  }

  pushChanges = toDoList => {
    const dto = toDoListToDto(toDoList);
    fetch("/api/to-do-list", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(dto)
    });
  }

  selectList = list => {
    return this.setState({
      activeToDoList: list,
      isListViewActive: false
    });
  }

  addList = () => {
    const newList = ToDoList.new("");
    this.state.toDoLists.push(newList);
    this.setState({
      activeToDoList: newList,
      toDoLists: this.state.toDoLists,
      isListViewActive: false
    });
  }

  activateListView = () => {
    this.setState({ isListViewActive: true });
  }

  renameList = title => {
    this.state.activeToDoList.rename(title);
    this.setState({ activeToDoList: this.state.activeToDoList });
    this.pushChanges(this.state.activeToDoList);
  }

  checkItem = item => {
    if (item.checked) {
      this.state.activeToDoList.uncheckItem(item.id);
    } else {
      this.state.activeToDoList.checkItem(item.id);
    }
    this.setState({ activeToDoList: this.state.activeToDoList });
    this.pushChanges(this.state.activeToDoList);
  }

  addItem = title => {
    this.state.activeToDoList.addItem(title);
    this.setState({ activeToDoList: this.state.activeToDoList });
    this.pushChanges(this.state.activeToDoList);
  }

  removeItem = item => {
    this.state.activeToDoList.removeItem(item.id);
    this.setState({ activeToDoList: this.state.activeToDoList });
    this.pushChanges(this.state.activeToDoList);
  }

  moveItem = (id, targetIndex) => {
    this.state.activeToDoList.moveItem(id, targetIndex);
    this.setState({ activeToDoList: this.state.activeToDoList });
    this.pushChanges(this.state.activeToDoList);
  }

  renameItem = (item, title) => {
    this.state.activeToDoList.renameItem(item.id, title);
    this.setState({ activeToDoList: this.state.activeToDoList });
    this.pushChanges(this.state.activeToDoList);
  }

  render() {
    return (
      <div className={"App" + (this.state.isListViewActive ? " overview" : "")}>
        <DetailView
          toDoList={this.state.activeToDoList}
          checkItem={this.checkItem}
          addItem={this.addItem}
          removeItem={this.removeItem}
          moveItem={this.moveItem}
          renameItem={this.renameItem}
          renameList={this.renameList}
          activateListView={this.activateListView}
        />

        <ListView
          toDoLists={this.state.toDoLists}
          selectList={this.selectList}
          addList={this.addList}
        />
      </div>
    );
  }
}



export default App;
