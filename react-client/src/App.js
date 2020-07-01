import React from 'react';

import './App.css';
import DetailView from './solvent/render/DetailView'
import ListView from './solvent/render/ListView'

import { notebookFromDto, notebookToDto } from './solvent/Dto'

class App extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      notebook: null,
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

  updateNotebook = f => {
    const notebook = this.state.notebook;
    const newNotebook = f(notebook);
    this.setState({ notebook: newNotebook });

    const activeToDoList = this.state.activeToDoList;
    if (activeToDoList) {
      const newActiveToDoList = notebook.getList(activeToDoList.id);
      if (newActiveToDoList) {
        const mergedActiveToDoList = activeToDoList.merge(newActiveToDoList);
        this.setState({ activeToDoList: mergedActiveToDoList });
      }
    }
  }

  updateActiveList = f => {
    this.updateNotebook(notebook => {
      const activeList = notebook.getList(this.state.activeToDoList.id);
      f(activeList);
      return notebook;
    })
  }

  syncState = async () => {
    if (this.state.notebook) {
      const newNotebook = await this.pushState(this.state.notebook);
      this.updateNotebook(notebook => notebook.merge(newNotebook));
    } else {
      const newNotebook = await this.fetchState();
      this.setState({ notebook: newNotebook });
    }
  }

  pushState = async notebook => {
    const dto = notebookToDto(notebook);
    const response = await fetch("api/notebook", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(dto)
    });
    const responseBody = await response.json();
    return notebookFromDto(responseBody);
  }

  fetchState = async () => {
    // TODO: Fetch for real user
    const response = await fetch("api/notebook/00000000-0000-0000-0000-000000000000");
    const responseBody = await response.json();
    return notebookFromDto(responseBody);
  }

  selectList = list => {
    return this.setState({
      activeToDoList: list,
      isListViewActive: false
    });
  }

  addList = () => {
    const list = this.state.notebook.addList("")
    this.setState({
      activeToDoList: list,
      notebook: this.state.notebook,
      isListViewActive: false
    });
  }

  activateListView = () => {
    this.setState({ isListViewActive: true });
  }

  renameList = title => {
    this.updateActiveList(list => list.rename(title));
  }

  checkItem = item => {
    if (item.checked) {
      this.updateActiveList(list => list.uncheckItem(item.id));
    } else {
      this.updateActiveList(list => list.checkItem(item.id));
    }
  }

  addItem = title => {
    this.updateActiveList(list => list.addItem(title));
  }

  removeItem = item => {
    this.updateActiveList(list => list.removeItem(item.id));
  }

  moveItem = (id, targetIndex) => {
    this.updateActiveList(list => list.moveItem(id, targetIndex));
  }

  renameItem = (item, title) => {
    this.updateActiveList(list => list.renameItem(item.id, title));
  }

  render() {
    return (
      <div className={"App" + (this.state.isListViewActive ? " overview" : "")}>
        <div className="ViewContainer">
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
        </div>

        <div className="ViewContainer">
          <ListView
            toDoLists={this.state.notebook ? this.state.notebook.getLists() : []}
            selectList={this.selectList}
            addList={this.addList}
          />
        </div>
      </div>
    );
  }
}

export default App;
