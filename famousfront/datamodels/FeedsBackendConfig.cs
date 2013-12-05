using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedsBackendConfig
  {
    [DataMember(EmitDefaultValue = false)]
    public string ip { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint port { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string db_address { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string db_name { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string[] categories { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string data_dir { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public ulong usage { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string image { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string thumbnail { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string document { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string feed_source { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string feed_entry { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public string proxy { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint summary_threshold { get; set; }
    [DataMember(EmitDefaultValue = false)]
    public uint thumbnail_width { get; set; }
  }
}