class ToDoItem {

  constructor(id, title, checked, orderValue, updatedAt) {
    this.id = id;
    this.title = title;
    this.checked = checked;
    this.orderValue = orderValue;
    this.updatedAt = updatedAt;
  }
}

export default ToDoItem;