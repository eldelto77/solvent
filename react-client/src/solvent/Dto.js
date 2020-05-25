import ToDoList from './ToDoList'
import ToDoItem from './ToDoItem'

function toDoItemFromDto(dto) {
  return new ToDoItem(
    dto.id,
    dto.title,
    dto.checked,
    dto.orderValue,
    dto.updatedAt
  );
}

function toDoItemToDto(item) {
  return {
    "id": item.id,
    "title": item.title,
    "checked": item.checked,
    "orderValue": item.orderValue,
    "updatedAt": item.updatedAt
  };
}

function toDoItemMapFromDto(dto) {
  const map = new Map();
  dto.map(toDoItemFromDto)
    .forEach(item => map.set(item.id, item));

  return map;
}

function toDoItemMapToDto(map) {
  const dto = [];
  map.forEach((item, _) => dto.push(toDoItemToDto(item)));

  return dto;
}

export function toDoListFromDto(dto) {
  return new ToDoList(
    dto.id,
    dto.title,
    toDoItemMapFromDto(dto.liveSet),
    toDoItemMapFromDto(dto.tombstoneSet),
    dto.updatedAt,
    dto.createdAt
  );
}

export function toDoListToDto(list) {
  return {
    "id": list.id,
    "title": list.title,
    "liveSet": toDoItemMapToDto(list.liveSet),
    "tombstoneSet": toDoItemMapToDto(list.tombstoneSet),
    "updatedAt": list.updatedAt,
    "createdAt": list.createdAt
  };
}