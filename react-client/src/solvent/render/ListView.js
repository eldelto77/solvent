import React from 'react';

import { ReactComponent as Magnify } from '../../icons/magnify.svg'
import { ReactComponent as PlusCircle } from '../../icons/plus-circle-shadow.svg'

export default function ListView(props) {
  return (
    <div className="ListView">
      <Header />
      <ListViewMain toDoLists={props.toDoLists} onClick={props.selectList} />
      <Footer onClick={props.addList} />
    </div>
  );
}

function Header(props) {
  return (
    <header>
      <div className="ListViewHeader header">
        <span className="HeaderSpacer"></span>
        <h1 className="HeaderTitle">Solvent</h1>
        <button className="DetailViewMenuButton menuButton" onClick={props.onSearchClick}>
          <Magnify />
        </button>
      </div>
    </header>
  );
}

function ListViewMain(props) {
  return (
    <div className="ListViewMain">
      <ToDoLists
        className="ListViewToDoListsOpen"
        title="Open"
        toDoLists={props.toDoLists.filter(list => !list.isChecked())}
        onClick={props.onClick}
      />

      <ToDoLists
        className="ListViewToDoListsDone"
        title="Done"
        toDoLists={props.toDoLists.filter(list => list.isChecked())}
        onClick={props.onClick}
      />
    </div>
  );
}

function ToDoLists(props) {
  return (
    <div>
      <span className="ListViewToDoListsTitle">{props.title}</span>
      {props.toDoLists.sort((a, b) => b.createdAt - a.createdAt)
        .map(toDoList =>
          <ToDoList key={toDoList.id} toDoList={toDoList} onClick={props.onClick} />
        )}
    </div>
  );
}

function ToDoList(props) {
  return (
    <button className={"ListViewToDoList" + (props.toDoList.isChecked() ? " checked" : "")} onClick={() => props.onClick(props.toDoList)}>
      <span className="ListViewToDoListTitle">{props.toDoList.title.value}</span>
    </button>
  );
}

function Footer(props) {
  return (
    <div className="ListViewFooter footer">
      <button className="ListViewAddButton" onClick={props.onClick}>
        <PlusCircle />
      </button>
    </div>
  );
}