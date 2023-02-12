function getData(){
    let name = document.getElementById('name').value
    let email = document.getElementById('email').value
    let phone = document.getElementById('phone').value
    let subject = document.getElementById('subject').value
    let message = document.getElementById('message').value

    if(name == "") {
        alert("Nama harus diisi")
    } else if(email == "") {
        alert("Email harus diisi")
    } else if(phone == "") {
        alert("Phone harus diisi")
    } else if(subject == "") {
        alert("Subject harus dipilih")
    } else if(message == "") {
        alert("Message harus diisi")
    }

    const defaultEmail = "udanieru@gmail.com"

    let mailTo = document.createElement('a')
    mailTo.href = `mailto:${defaultEmail}?subject=${subject}&body=Halo nama saya ${name}, saya ingin ${message} tolong hubungin saya ${phone}`
    mailTo.click()
}