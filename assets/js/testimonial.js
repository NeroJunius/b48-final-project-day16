
const promise = new Promise((resolve, reject) => {
    const gt = new XMLHttpRequest();
    gt.open("GET", "https://api.npoint.io/011aef2935141472432f", true);
    // console.log(xhr);
    gt.onload = () => {
        if (gt.status === 200) {
            resolve(JSON.parse(gt.response));
        } else {
            reject("Error loading card data.");
        };
    };

    gt.onerror = () => {
        reject("Network problem.");
    };

    gt.send();
});


async function showAll() {

    const response = await promise;

    let markUp = "";

    response.forEach(function (item) {
        markUp += `
        <div class="cards">
            <img src="${item.img}">
            <p class="quote">"${item.quote}"</p>
            <br>
            <br>
            <p class="user">from ${item.user}</p>
            <p class="star">${item.rating} <i class="fa-solid fa-star fa-fw"></i></p>
        </div>
        `;
    });

    document.getElementById("testimonial").innerHTML = markUp;
};

showAll();

async function sort(rating) {

    const response = await promise;

    let markUp = "";

    const sorted = response.filter(function (item) {
        return item.rating === rating;
    });

    if (sorted.length === 0) {
        markUp += `
        <h2>Empty</h2>
        `;
    } else {
        sorted.forEach(function (item) {
            markUp += `
            <div class="cards">
                <img src="${item.img}">
                <p class="quote">"${item.quote}"</p>
                <p class="user">${item.rating} <i class="fa-solid fa-star fa-fw"></i> from ${item.user}</p>
            </div>
            `;
        });
    };

    document.getElementById("testimonial").innerHTML = markUp;
};
