<template>
  <v-container>
    <v-layout
      text-center
      wrap
    >
      <v-flex xs12>
        <v-img
          :src="require('../assets/Acronis.svg')"
          class="my-3"
          contain
          height="200"
        ></v-img>
      </v-flex>
      <v-flex mb-4>
        <p class="subheading font-weight-regular" style="font-size:35px">
          Verify the authenticity of a document <br>signed by Acronis Notary
        </p>
      </v-flex>
      <v-flex xs12>
        <v-file-input
          v-model="files"
          color="deep-purple accent-4"
          counter
          label="File input"
          multiple
          placeholder="Select your files"
          prepend-icon="mdi-paperclip"
          outlined
          :show-size="1000"
        >
          <template v-slot:selection="{ index, text }">
            <v-chip
              v-if="index < 2"
              color="deep-purple accent-4"
              dark
              label
              small
            >
              {{ text }}
            </v-chip>

            <span
              v-else-if="index === 2"
              class="overline grey--text text--darken-3 mx-2"
            >
              +{{ files.length - 2 }} File(s)
            </span>
          </template>
        </v-file-input>
      </v-flex>
      <v-flex mb-4>
        <v-btn @click="verify"> Verify </v-btn>
      </v-flex>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
const axios = require("axios")
export default {
  data: () => ({
    files: [],
  }),
  methods: {
    verify() {
      console.log(this.files[0])
      axios.post("/api/uploadHash",{
        Name:this.files[0].name,
        LastModified: this.files[0].LastModified,
        Size: this.files[0].size,
        type: this.files[0].type
      }).then((response) => {
        console.log(response.data)
      })
    }
  }
};
</script>
