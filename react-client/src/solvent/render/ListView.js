import React from 'react';

import { ReactComponent as ArrowLeft } from '../../icons/arrow-left.svg'
import { ReactComponent as Magnify } from '../../icons/magnify.svg'
import { ReactComponent as PlusCircle } from '../../icons/plus-circle.svg'

export default function ListView(props) {
  return (
    <div className="ListView">
      <Header onClick={props.onBack} backButtonEnabled={props.backButtonEnabled} />
      <ListViewMain toDoLists={props.toDoLists} onClick={props.selectList} />
      <Footer onClick={props.addList} />
    </div>
  );
}

function Header(props) {
  return (
    <div className="ListViewHeader header">
      <span className="ListViewSearchBarLogo menuButton">
        <Magnify />
      </span>
      <input className="ListViewSearchBar" type="text" placeholder="Type to search" />
      <button className="ListViewBackButton menuButton" onClick={props.onClick} disabled={!props.backButtonEnabled}>
        <ArrowLeft />
      </button>
    </div>
  );
}

function ListViewMain(props) {
  return (
    <div className="ListViewMain">
      {props.toDoLists.sort((a, b) => b.createdAt - a.createdAt)
        .map(toDoList =>
          <ToDoList key={toDoList.id} toDoList={toDoList} onClick={props.onClick} />
        )}
    </div>
  );
}

function ToDoList(props) {
  return (
    <button className="ListViewToDoList" onClick={() => props.onClick(props.toDoList)}>
      {props.toDoList.title}
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