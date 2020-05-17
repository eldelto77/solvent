import ToDoItem from './ToDoItem'
import { v4 as uuid } from 'uuid'

class ToDoList {

  constructor(id, title, liveSet, tombstoneSet) {
    this.id = id;
    this.title = title;
    this.liveSet = liveSet;
    this.tombstoneSet = tombstoneSet;
  }

  static new(title) {
    return new ToDoList(uuid(), title, new Map(), new Map());
  }

  get items() {
    const items = [];
    this.liveView().forEach((item, _) => items.push(item))

    return items;
  }

  addItem(title) {
    const id = uuid();
    const item = new ToDoItem(id, title, false, this.nextOrderValue());
    this.liveSet.set(id, item);

    return id;
  }

  getItem(id) {
    return this.liveView().get(id);
  }

  removeItem(id) {
    const items = this.liveView();
    if (items.has(id)) {
      const item = items.get(id);
      this.tombstoneSet.set(id, item);
    }
  }

  checkItem(id) {
      const item = this.getItem(id);
      item.checked = true

      return id;
  }

  uncheckItem(id) {
      const item = this.getItem(id);
      this.tombstoneSet.set(id, item);
      const newId = uuid();
      const newItem = new ToDoItem(newId, item.title, false, item.orderValue);
      this.liveSet.set(newId, newItem)

      return newId;
  }

  moveItem(id, targetIndex) {
      const item = this.getItem(id);
      const items = this.items.sort((a, b) => a.orderValue - b.orderValue);

      const orderValueMid = items[targetIndex].orderValue;
      let orderValueAdjacent;
      if (orderValueMid < item.orderValue) {
        // Movint item up
        if ((targetIndex - 1) >= 0) {
          orderValueAdjacent = items[targetIndex - 1].orderValue;
        } else {
          orderValueAdjacent = 0.0;
        }
      } else if (orderValueMid > item.orderValue) {
        // Movint item down
        if ((targetIndex + 1) < items.length) {
          orderValueAdjacent = items[targetIndex + 1].orderValue;
        } else {
          orderValueAdjacent = this.nextOrderValue();
        }
      } else {
        return;
      }

      const newOrderValue = (orderValueMid + orderValueAdjacent) / 2.0;
      item.orderValue = newOrderValue;
      this.liveSet.set(item.id, item);

      return item.id;
  }

  renameItem(id, title) {
    const oldItem = this.getItem(id);
    this.removeItem(oldItem.id);
    const newId = this.addItem(title);
    const newItem = this.getItem(newId);
    newItem.orderValue = oldItem.orderValue;
    this.liveSet.set(newItem.id, newItem);

    return newId;
  }

  liveView() {
    const liveView = new Map();
    this.liveSet.forEach((item, id) => {
      if (!this.tombstoneSet.has(id)) {
        liveView.set(id, item);
      }
    })

    return liveView;
  }

  nextOrderValue() {
    let orderValue = 0.0;
    this.liveView().forEach((item, _) => {
      if (item.orderValue > orderValue) {
        orderValue = item.orderValue;
      }
    });

    return orderValue + 10.0;
  }
}

export default ToDoList;