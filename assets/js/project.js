//checkboxes
function getCheckedCheckboxesFor() {
    var checkboxes = document.querySelectorAll('input[name="technologies"]:checked'),
        values = [];
    Array.prototype.forEach.call(checkboxes, function (el) {
        values.push(el.value);
    });
    return values;
}

let projects = [];

function getData(event) {
    event.preventDefault();

    let title = document.getElementById("project_name").value;
    let startDate = new Date(document.getElementById("start_date").value);
    let endDate = new Date(document.getElementById("end_date").value);
    let description = document.getElementById("description").value;
    let image = document.getElementById("image_uploads").files;

    image = URL.createObjectURL(image[0]);

    let checked = getCheckedCheckboxesFor();

    let dataProject = {
        title,
        startDate,
        endDate,
        description,
        checked,
        image,
    };

    projects.push(dataProject);
    showData();
}

function getDuration(startDate, endDate) {
    let duration = endDate - startDate;

    //pembulatan
    let durationMonth = Math.floor(duration / 1000 / 60 / 60 / 24 / 30);
    let durationDay = Math.floor(duration / 1000 / 60 / 60 / 24);

    if (durationMonth > 0) {
        return durationMonth + " Bulan";
    } else if (durationDay > 0) {
        return durationDay + " Hari";
    }
}

function showData() {
    document.getElementById("content-project").innerHTML = "";
    for (let i = 0; i <= projects.length; i++) {
        document.getElementById("content-project").innerHTML += `
            <a href="detail-project.html">
                <div class="card-project">
                    <div class="img-container">
                        <img src="${projects[i].image}" alt="" width="250" style="border-radius:10px">
                    </div>
                    <div>
                        <h2 style="margin: 10px 0px 5px 0px; color: black;">${projects[i].title}</h2>
                        <p style=" margin-bottom: 10px; color: gray;">Durasi : ${getDuration(projects[i].startDate, projects[i].endDate)}</p>
                    </div>
                    <p style="color: black;">${projects[i].description}</p>
                    <div style="margin-top: 15px;">
                        <span style="color: black;">${projects[i].checked.join("&nbsp&nbsp&nbsp&nbsp")}</span>
                        <div style="display: flex; justify-content: space-between; margin-top:10px; column-gap: 10px">
                            <button class="btn-primary">Edit</button>
                            <button class="btn-danger">Delete</button>
                        </div>
                    </div>
                </div>
            </a>
        `;
    }
}
