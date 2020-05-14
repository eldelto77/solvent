import React from 'react'

class RToDoList extends React.Component {

  constructor(props) {
    super(props);

    this.checkItem = this.checkItem.bind(this);
  }

  checkItem(item) {
    return this.props.checkItem(item);
  }

  render() {
    return(
    <div className="ToDoList">
      <h1>{this.props.toDoList.title}</h1>

      {this.props.toDoList.items.map(item => {
        return (
        <RToDoItem key={item.id} item={item} onCheck={() => this.checkItem(item)} />
        );
      })}
    </div>
    );
  }
}

  function RToDoItem(props) {
    return (
      <div>
        <input 
        type="checkbox" 
        checked={props.item.checked} 
        onChange={props.onCheck}
        />
        {props.item.title}
      </div>
    );
  }

export default RToDoList;