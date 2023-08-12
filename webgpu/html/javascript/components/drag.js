export class Drag {
  constructor (overlay, getXY, setXY) {
    this.overlay = overlay
    this.getXY = getXY
    this.setXY = setXY

    this.dragging = false
    this.hscale = 2 // FIXME calculate from client width
    this.vscale = 2 // FIXME calculate from client height

    this.origin = { x: 0, y: 0 }
    this.startXY = { x: 0, y: 0 }
  }

  start (event, canvas) {
    this.dragging = true
    this.origin = this.getXY()
    this.startXY = { x: this.hscale * event.offsetX, y: this.vscale * event.offsetY }

    canvas.onpointermove = (event) => this.onPointerMove(event, canvas)
    canvas.onpointerup = (event) => this.onPointerUp(event, canvas)

    canvas.setPointerCapture(event.pointerId)
  }

  onPointerMove (event, canvas) {
    if (this.dragging) {
      const xy = {
        x: this.hscale * event.offsetX,
        y: this.vscale * event.offsetY
      }

      const dx = xy.x - this.startXY.x
      const dy = xy.y - this.startXY.y

      this.setXY(this.origin.x + dx, this.origin.y + dy, this.dragging)
    }
  }

  onPointerUp (event, canvas, drag) {
    canvas.onpointermove = null
    canvas.onpointerup = null
    canvas.releasePointerCapture(event.pointerId)

    if (this.dragging) {
      const xy = { x: this.hscale * event.offsetX, y: this.vscale * event.offsetY }
      const dx = xy.x - this.startXY.x
      const dy = xy.y - this.startXY.y

      this.dragging = false
      this.setXY(this.origin.x + dx, this.origin.y + dy, this.dragging)
    }
  }
}
