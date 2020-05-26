import React from 'react';

import RToDoList from './RToDoList'

export default function DetailView(props) {
  return (
    <div className="DetailView">
      <Header onMenuClick={props.activateListView} />
      {props.toDoList ?
        <RToDoList
          toDoList={props.toDoList}
          checkItem={props.checkItem}
          addItem={props.addItem}
          removeItem={props.removeItem}
          moveItem={props.moveItem}
          renameItem={props.renameItem}
          renameList={props.renameList}
        />
        : ""}
    </div>
  );
}

function Header(props) {
  return (
    <header>
      <div className="DetailViewHeader header">
        <span className="HeaderSpacer"></span>
        <h1 className="HeaderTitle">Solvent</h1>
        <button className="DetailViewMenuButton menuButton" onClick={props.onMenuClick}></button>
      </div>
    </header>
  );
}