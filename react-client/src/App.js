import React from 'react';

import './App.css';
import DetailView from './solvent/render/DetailView'
import ListView from './solvent/render/ListView'

import ToDoList from './solvent/ToDoList'

import { toDoListFromDto, toDoListToDto } from './solvent/Dto'

class App extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      toDoLists: [],
      activeToDoList: null,
      isListViewActive: true
    }
  }

  componentDidMount() {
    this.syncState();
    this.timer = setInterval(() => this.syncState(), 1000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
    this.timer = null;
  }

  syncState = async () => {
    await this.pushState(this.state.toDoLists);
    const newToDoLists = await this.fetchState();

    const oldToDoListMap = new Map();
    this.state.toDoLists.forEach(list => oldToDoListMap.set(list.id, list));

    const mergedToDoLists = [];
    newToDoLists.forEach(newToDoList => {
      if (oldToDoListMap.has(newToDoList.id)) {
        const oldToDoList = oldToDoListMap.get(newToDoList.id);
        const mergedToDoList = oldToDoList.merge(newToDoList);
        mergedToDoLists.push(mergedToDoList);
      } else {
        mergedToDoLists.push(newToDoList);
      }
    });

    this.setState({ toDoLists: mergedToDoLists });

    if (this.state.activeToDoList) {
      const activeToDoList = mergedToDoLists.find(list => list.id === this.state.activeToDoList.id);
      if (activeToDoList) {
        this.setState({ activeToDoList: activeToDoList });
      }
    }
  }

  pushState = async toDoLists => {
    const dtos = toDoLists.map(toDoListToDto);
    const requestBody = {
      toDoLists: dtos
    }

    return fetch("api/to-do-list/bulk", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(requestBody)
    });
  }

  fetchState = async () => {
    const response = await fetch("api/to-do-list");
    const responseBody = await response.json();
    return responseBody.toDoLists.map(toDoListFromDto);
  }

  backToDetailView = () => {
    return this.setState({ isListViewActive: false });
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
  }

  checkItem = item => {
    if (item.checked) {
      this.state.activeToDoList.uncheckItem(item.id);
    } else {
      this.state.activeToDoList.checkItem(item.id);
    }
    this.setState({ activeToDoList: this.state.activeToDoList });
  }

  addItem = title => {
    this.state.activeToDoList.addItem(title);
    this.setState({ activeToDoList: this.state.activeToDoList });
  }

  removeItem = item => {
    this.state.activeToDoList.removeItem(item.id);
    this.setState({ activeToDoList: this.state.activeToDoList });
  }

  moveItem = (id, targetIndex) => {
    this.state.activeToDoList.moveItem(id, targetIndex);
    this.setState({ activeToDoList: this.state.activeToDoList });
  }

  renameItem = (item, title) => {
    this.state.activeToDoList.renameItem(item.id, title);
    this.setState({ activeToDoList: this.state.activeToDoList });
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
          onBack={this.backToDetailView}
          backButtonEnabled={this.state.activeToDoList ? true : false}
        />
      </div>
    );
  }
}

export default App;
