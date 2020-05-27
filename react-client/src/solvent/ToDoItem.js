class ToDoItem {

  constructor(id, title, checked, orderValue, updatedAt) {
    this.id = id;
    this.title = title;
    this.checked = checked;
    this.orderValue = orderValue;
    this.updatedAt = updatedAt;
  }

  merge(other) {
    const checked = this.checked || other.checked

    let orderValue = this.orderValue;
    let updatedAt = this.updatedAt;
    if (other.updatedAt > this.updatedAt) {
      orderValue = other.orderValue;
      updatedAt = other.updatedAt;
    }

    return new ToDoItem(this.id, this.title, checked, orderValue, updatedAt);
  }
}

export default ToDoItem;