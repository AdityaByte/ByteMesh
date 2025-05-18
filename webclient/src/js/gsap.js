document.addEventListener("DOMContentLoaded", function(event) {

    // So what we have to do here we have to create a timeline for the GSAP
    const headerElements = document.querySelector("header")
    const sectionElements = document.querySelector(".section1")
    if (headerElements) {
        timeline = gsap.timeline({paused: true})
        timeline
            .add(navbarAnimation())
            .add(headerAnimation())
            .add(loadCardsAndAnimate())
            timeline.play()
    }
})

const cardAnimation = () => {
    loadCardsAndAnimate()
}

const headerAnimation = () => {

    tl = gsap.timeline()

    tl.from("header h1", {
        y: -100,
        opacity: 0,
        duration: 0.3,
    })

    tl.from("header p", {
        y: -100,
        opacity: 0,
        duration: 0.3,
    })

    tl.from("header button", {
        opacity: 0,
        duration: 0.3,
    })

    return tl
}

const navbarAnimation = () => {
    const navlogo = document.querySelector(".logo")
    const ulItems = document.querySelectorAll("nav ul a")

    let timeline = gsap.timeline()

    timeline.from(navlogo, {
        y: -10,
        opacity: 0,
        duration: 0.3,
    })

    timeline.from(ulItems, {
        y: -10,
        opacity: 0,
        duration: 0.3,
        stagger: 0.2,
    })

    return timeline
}


const loadCardsAndAnimate = function() {
    // Here what we have to set a timeout ok dude for sometime so the page has been fully loaded thereafter we need to create the 3 child ok.
    setTimeout(function(){
        // Here we have to fetch out the parent section.
        const section1 = document.querySelector(".section1")
        // Now we have to select the proper datatype ok.
        const myData = [
            {
                title: "Fault Tolerant",
                imgSrc: "https://imgs.search.brave.com/Ni5pT-GLiVXuLpUFGRPXoE0p-r5fpC5NpC7YjmCOkCY/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly90My5m/dGNkbi5uZXQvanBn/LzEwLzc5LzI1LzU4/LzM2MF9GXzEwNzky/NTU4ODFfZ080a0Rp/cDVpNGwzUWJ2Qktj/U0VyWktXR0swTllx/ZGguanBn",
                content: "The system continues operating even when one or more components fail, ensuring high availability."
            },
            {
                title: "Scalable Architecture",
                imgSrc: "https://imgs.search.brave.com/Ni5pT-GLiVXuLpUFGRPXoE0p-r5fpC5NpC7YjmCOkCY/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly90My5m/dGNkbi5uZXQvanBn/LzEwLzc5LzI1LzU4/LzM2MF9GXzEwNzky/NTU4ODFfZ080a0Rp/cDVpNGwzUWJ2Qktj/U0VyWktXR0swTllx/ZGguanBn",
                content: "Designed to handle growing amounts of work and data by adding more nodes without downtime."
            },
            {
                title: "High Throughput",
                imgSrc: "https://imgs.search.brave.com/Ni5pT-GLiVXuLpUFGRPXoE0p-r5fpC5NpC7YjmCOkCY/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly90My5m/dGNkbi5uZXQvanBn/LzEwLzc5LzI1LzU4/LzM2MF9GXzEwNzky/NTU4ODFfZ080a0Rp/cDVpNGwzUWJ2Qktj/U0VyWktXR0swTllx/ZGguanBn",
                content: "Supports fast read/write operations across distributed nodes, making it ideal for big data processing."
            }
        ];


        // here we need to create the three child of it and append the child to the section and make it possible to load out.
        for(let i=0; i<3; i++) {


            // Here we need to create the card first
            const card = document.createElement("div")
            card.className = "card"

            // Here we have to create a new child
            const titleChild = document.createElement("span")
            titleChild.innerHTML = myData[i].title
            const imageChild = document.createElement("img")
            imageChild.style.width ="50px"
            imageChild.style.height = "50px"
            imageChild.src = myData[i].imgSrc
            const contentChild = document.createElement("span")
            contentChild.innerHTML = myData[i].content

            card.append(titleChild, imageChild, contentChild)

            section1.appendChild(card)
            console.log("card data has been set.")
        }

        gsap.from(".section1 .card", {
            y: -20,
            opacity: 0,
            duration: 0.4,
            stagger: 0.2,
            scrollTrigger: {
                trigger: ".section1",
                scroller: "body",
                start: "top 50%",
            }
        });

        console.log("Cards loaded and animated.");
        return tl

    }, 10) // Delays for 1 second.

}