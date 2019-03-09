# TODO: при нажатии на сам файл выходит загрузка, а не предупреждение!
pop = new Vue
cont = new Vue
  el: '#container'
  data:
    endpoint: 'http://localhost:8080/get'
    deleteEndpoint: 'http://localhost:8080/delete'
    uploadEndpoint: 'http://localhost:8080/upload'
    dirs: []
    files: [
      {
        Type: ''
        Text: ''
        Checked: false
      }
    ]
    dirName: ''
    empty: true
    baseChecked: false
    deletePopup: false
    uploadPopup: false
  methods:
    getAllDirs: () ->
      this.$http.get(this.endpoint).then((response) ->
        this.dirs = response.body
        (err) ->
          console.log 'Error with handling data: ', err
      )
    getFiles: (path, type) ->
      if type != 'md-folder'
        alert "It is file!!!"
        return
      this.dirName = path + "/"
      if path == ''
        this.dirName = ''
        this.empty = true
        this.files = []
        return

      options =
        params:
          path: path

      this.$http.get(this.endpoint, options).then((response) ->
          this.files = []
          this.files = response.body
          console.log(this.files)
          this.empty = this.files == [] || this.files.length == 0
          console.log(this.empty)
        (err) ->
          console.log 'Error with handling data: ', err
      )

    filterUp: (value) ->
      str = value.split("/")
      value = ''
      if str.length <= 2
        return ''
      else
        for i in [0...str.length - 3]
          value += str[i] + "/"
      value += str[str.length - 3]
      return value

    checkAll: () ->
      this.baseChecked = !this.baseChecked
      console.log(this.baseChecked)
      for f in this.files
        f.Checked = this.baseChecked

    enableDeletePopupIfNotEmpty: () ->
      isNotEmpty = false
      for f in this.files
        if f.Checked
          isNotEmpty = true
          break
      this.deletePopup = isNotEmpty
  created: () ->
    this.getAllDirs()

pop = new Vue
  el: '#popups'
  data:
    newProjectPopup: false
    projectName: ''
    createEndpoint: 'http://localhost:8080/create'
  methods:
    deleteFiles: () ->
      cont.deletePopup = false
      data1 = []
      for f in cont.files
        if f.Checked
          data1.push(cont.dirName + f.Name)
      data = JSON.stringify({
          FileNames: data1
        }
      )
      this.$http.post(cont.deleteEndpoint, data).then((result) ->
        cont.getFiles(cont.dirName.substring(0, cont.dirName.length - 1), "md-folder")
        (err) ->
          console.log("Error: ", err)
      )

    changeDeletePopup: () ->
      cont.deletePopup = !cont.deletePopup

    getDeletePopup: () ->
      console.log(cont.deletePopup)
      return cont.deletePopup

    getUploadPopup: () ->
      console.log(cont.uploadPopup)
      return cont.uploadPopup

    changeUploadPopup: () ->
      cont.uploadPopup = !cont.uploadPopup

    changeNewPopup: () ->
      this.newProjectPopup = !this.newProjectPopup

    uploadFiles: () ->
      filesToUpload = new FormData()
      fp = document.getElementById('file-field').files
      console.log("A")
      for e in fp
        console.log(e)
        filesToUpload.append("Files", e)
      #data = JSON.stringify({
       #   Directory: cont.dirName
        #  Files: filesToUpload
        #}
      #)
      filesToUpload.append("Directory", cont.dirName)
      this.$http.post(cont.uploadEndpoint, filesToUpload).then((result) ->
          console.log(result)
        (err) ->
          console.log("Error: ", err)
      )
      this.changeUploadPopup()
      cont.getFiles(cont.dirName.substring(0, cont.dirName.length - 1), "md-folder")

    createProject: () ->
      createForm = new FormData()
      createForm.append("Name", this.projectName)
      console.log createForm
      this.$http.post(this.createEndpoint, createForm).then((result) ->
          console.log(result)
        (err) ->
          console.log("Error: ", err)
      )
      this.changeNewPopup()
      cont.getAllDirs()

nav = new Vue
  el: '#nav'
  methods:
    showNewProject: () ->
      pop.newProjectPopup = true
      console.log(".....", pop.newProjectPopup)